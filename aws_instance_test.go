package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"testing"
	"time"
)

func TestIsInLatestScan(t *testing.T) {
	testTime := "19 Nov 18 00:00 PDT"
	Now = func() time.Time {
		duration, _ := time.ParseDuration("2h")
		parsed, _ := time.Parse(time.RFC822, testTime)
		return parsed.Add(duration)
	}

	scanInterval := &ec2.Tag{Key: aws.String("ScanInterval"), Value: aws.String("72h")}
	lastScanned := &ec2.Tag{Key: aws.String("LastScanned"), Value: aws.String(testTime)}
	rawInstance := &ec2.Instance{Tags: []*ec2.Tag{scanInterval, lastScanned}}
	instance := &AwsInstance{rawInstance: rawInstance, config: NewConfig()}
	instance.LoadTags()
	if res := instance.IsInLatestScan(); res != "Yes" {
		t.Errorf("IsInLatestScan should have been %s but was %s", "Yes", res)
	}
}
