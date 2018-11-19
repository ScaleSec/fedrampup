package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Fetcher struct {
	config    Config
	instances []*Instance
	sess      *session.Session
}

func NewFetcher(config Config, sess *session.Session) *Fetcher {
	return &Fetcher{
		config:    config,
		sess:      sess,
		instances: []*Instance{},
	}
}

func (this *Fetcher) Run() ([]*Instance, error) {
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
	return this.instances, err
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
				this.instances = append(this.instances, &Instance{
					rawInstance: instance,
					rawImage:    image,
					config:      this.config,
				})
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
