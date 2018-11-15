package fedrampup

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
)

type Fetcher struct {
  config Config
  instances []*Instance
}

func NewFetcher(config Config) (*Fetcher){
  return &Fetcher{
    config: config,
    instances: []*Instance,
  }
}

func (this *Fetcher) Run() ([]Instance, error) {
  var err error
  if len(this.config.Profiles) > 0 {
    for profile := range this.config.Profiles {
      for region := range this.config.Regions {
        sess := session.Must(session.NewSession(&aws.Config{
          Credentials: credentials.NewSharedCredentials("", profile),
          Region: aws.String(region),
        }))
        client := ec2.New(sess)
        err = RunBatch(client)
      }
    }
  }  else {
    for region := range this.config.Regions {
      sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String(region),
      }))
      client := ec2.New(sess)
      err = RunBatch(client)
    }
  }
  return this.instances, err
}

func (this *Fetcher) RunBatch(client ec2.EC2) error {
  instancesResult, err := client.DescribeInstances(nil)
  if err != nil {
    return err
  }

  for {
    for reservation := range instancesResult.Reservations {
      for instance := range reservation.Instances {
        imagesResult, err := svc.DescribeImages(&ec2.DescribeImagesInput{
          ImageIds: []*string{ instance.ImageId },
        }).Images[0]
        append(this.instances, &Instance{
          rawInstance: instance,
          rawImage: image,
          config: this.config,
        })
      }
    }

    // Pagination
    if result.NextToken {
      result, err := svc.DescribeInstances(&DescribeInstancesInput{
        NextToken: result.NextToken,
      })
    } else { break }
  }
}
