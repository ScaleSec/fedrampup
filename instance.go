package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Instance struct {
	rawImage    *ec2.Image
	rawInstance *ec2.Instance
	config      Config
	tags        map[string]string
}

func (this *Instance) Row() []string {
	t := this.config.Tags
	inst := this.rawInstance
	img := this.rawImage
	this.LoadTags()

	return []string{
		aws.StringValue(inst.InstanceId),
		aws.StringValue(inst.PrivateIpAddress),
		this.IsVirtual(),
		this.IsPublic(),
		aws.StringValue(inst.PrivateDnsName),
		this.tags[t.netbios],
		this.MacAddress(),
		this.tags[t.authenticatedScanPlanned],
		this.tags[t.baselineConfig],
		aws.StringValue(img.Description),
		aws.StringValue(inst.Placement.AvailabilityZone),
		this.tags[t.assetType],
		aws.StringValue(inst.InstanceType),
		this.IsInLatestScan(),
		this.tags[t.applicationVendor],
		this.tags[t.applicationNameAndVersion],
		this.tags[t.applicationPatchLevel],
		this.tags[t.applicationFunction],
		this.tags[t.comments],
		this.tags[t.serialNumber],
		aws.StringValue(inst.SubnetId),
		this.tags[t.sysadmin],
		this.tags[t.appadmin],
	}
}

func (this *Instance) IsVirtual() string {
	if aws.StringValue(this.rawInstance.VirtualizationType) == "HVM" {
		return "Yes"
	} else {
		return "No"
	}
}

func (this *Instance) IsPublic() string {
	if this.rawInstance.PublicIpAddress != nil {
		return "Yes"
	} else {
		return "No"
	}
}

func (this *Instance) MacAddress() string {
	return aws.StringValue(this.rawInstance.NetworkInterfaces[0].MacAddress)
}

// NOTE: We expect there to be two tags on the instance to determine
// if it was in the latest scan or not. We need the interval of the scan
// such as 1d or 1w, etc. and then we expect the lastScanned tag to be in
// RFC822 format.
// TODO: Maybe the user would want to pass in their own format.
func (this *Instance) IsInLatestScan() string {
	interval, err := time.ParseDuration(this.config.ScanInterval)
	if err != nil {
		return ""
	}
	// Format: "02 Jan 06 15:04 MST"
	// https://www.w3.org/Protocols/rfc822/ "Date and Time Specification"
	lastScanned, err := time.Parse(time.RFC822, this.tags[this.config.Tags.lastScanned])
	if err != nil {
		return ""
	}

	// NOTE: time.Sub() gives a duration, but time.Add()
	// returns a time, which we need, so we'll add the negative
	if lastScanned.After(time.Now().Add(interval * -1)) {
		return "Yes"
	} else {
		return "No"
	}
}

func (this *Instance) LoadTags() {
	this.tags = make(map[string]string)
	for _, tag := range this.rawInstance.Tags {
		this.tags[aws.StringValue(tag.Key)] = aws.StringValue(tag.Value)
	}
}
