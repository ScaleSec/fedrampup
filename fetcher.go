package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Fetcher struct {
	config    *Config
	instances []*AwsInstance
	sess      *session.Session
	journaler *Journaler
}

func NewFetcher(config *Config, sess *session.Session, journaler *Journaler) *Fetcher {
	return &Fetcher{
		config:    config,
		sess:      sess,
		instances: []*AwsInstance{},
		journaler: journaler,
	}
}

func (this *Fetcher) Run() (string, error) {
	var err error
	if len(this.config.Roles) > 0 {
		for _, role := range this.config.Roles {
			for _, region := range this.config.Regions {
				creds := stscreds.NewCredentials(this.sess, role)
				client := ec2.New(this.sess, aws.NewConfig().WithCredentials(creds).WithRegion(region))
				err = this.RunBatch(client)
			}
		}
	} else {
		for _, region := range this.config.Regions {
			client := ec2.New(this.sess, aws.NewConfig().WithRegion(region))
			err = this.RunBatch(client)
		}
	}
	return this.journaler.Filename, err
}

func (this *Fetcher) RunBatch(client *ec2.EC2) error {
	instancesResult, err := client.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		return err
	}

	for {
		for _, reservation := range instancesResult.Reservations {
			for _, instance := range reservation.Instances {
				imagesResult, err := client.DescribeImages(&ec2.DescribeImagesInput{
					ImageIds: []*string{instance.ImageId},
				})
				if err != nil {
					return err
				}
				image := imagesResult.Images[0]
				instance := NewAwsInstance(instance, image, this.config)

				// Add journal entry
				this.journaler.Write(instance.Row())
			}
		}

		// Pagination
		if instancesResult.NextToken != nil {
			instancesResult, err = client.DescribeInstances(&ec2.DescribeInstancesInput{
				NextToken: instancesResult.NextToken,
			})
		} else {
			break
		}
	}
	return err
}
