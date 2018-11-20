package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"time"
)

//Gives the ability to stub out Now in tests
var Now = time.Now

type AwsInstance struct {
	rawImage                  *ec2.Image
	rawInstance               *ec2.Instance
	config                    *Config
	tags                      map[string]string
	Id                        string
	IpAddress                 string
	Virtual                   string
	Public                    string
	DNS                       string
	NetBIOS                   string
	MacAddress                string
	AuthenticatedScanPlanned  string
	BaselineConfiguration     string
	OS                        string
	Location                  string
	AssetType                 string
	Hardware                  string
	InLatestScan              string
	ApplicationVendor         string
	ApplicationNameAndVersion string
	ApplicationPatchLevel     string
	ApplicationFunction       string
	Comments                  string
	SerialNumber              string
	Network                   string
	SysAdmin                  string
	AppAdmin                  string
}

func NewAwsInstance(inst *ec2.Instance, img *ec2.Image, config *Config) *AwsInstance {
	instance := &AwsInstance{
		rawInstance: inst,
		rawImage:    img,
		config:      config,
	}
	t := instance.config.Tags
	instance.LoadTags()

	instance.Id = aws.StringValue(inst.InstanceId)
	instance.IpAddress = aws.StringValue(inst.PrivateIpAddress)
	instance.Virtual = instance.IsVirtual()
	instance.Public = instance.IsPublic()
	instance.DNS = aws.StringValue(inst.PrivateDnsName)
	instance.NetBIOS = instance.tags[t.netbios]
	instance.MacAddress = instance.GetMacAddress()
	instance.AuthenticatedScanPlanned = instance.tags[t.authenticatedScanPlanned]
	instance.BaselineConfiguration = instance.tags[t.baselineConfig]
	instance.OS = aws.StringValue(img.Description)
	instance.Location = aws.StringValue(inst.Placement.AvailabilityZone)
	instance.AssetType = instance.tags[t.assetType]
	instance.Hardware = aws.StringValue(inst.InstanceType)
	instance.InLatestScan = instance.IsInLatestScan()
	instance.ApplicationVendor = instance.tags[t.applicationVendor]
	instance.ApplicationNameAndVersion = instance.tags[t.applicationNameAndVersion]
	instance.ApplicationPatchLevel = instance.tags[t.applicationPatchLevel]
	instance.ApplicationFunction = instance.tags[t.applicationFunction]
	instance.Comments = instance.tags[t.comments]
	instance.SerialNumber = instance.tags[t.serialNumber]
	instance.Network = aws.StringValue(inst.SubnetId)
	instance.SysAdmin = instance.tags[t.sysadmin]
	instance.AppAdmin = instance.tags[t.appadmin]
	return instance
}

func (this *AwsInstance) Row() []string {
	return []string{
		this.Id,
		this.IpAddress,
		this.Virtual,
		this.Public,
		this.DNS,
		this.NetBIOS,
		this.MacAddress,
		this.AuthenticatedScanPlanned,
		this.BaselineConfiguration,
		this.OS,
		this.Location,
		this.AssetType,
		this.Hardware,
		this.InLatestScan,
		this.ApplicationVendor,
		this.ApplicationNameAndVersion,
		this.ApplicationPatchLevel,
		this.ApplicationFunction,
		this.Comments,
		this.SerialNumber,
		this.Network,
		this.SysAdmin,
		this.AppAdmin,
	}
}

func (this *AwsInstance) IsVirtual() string {
	if aws.StringValue(this.rawInstance.VirtualizationType) == "HVM" {
		return "Yes"
	} else {
		return "No"
	}
}

func (this *AwsInstance) IsPublic() string {
	if this.rawInstance.PublicIpAddress != nil {
		return "Yes"
	} else {
		return "No"
	}
}

func (this *AwsInstance) GetMacAddress() string {
	return aws.StringValue(this.rawInstance.NetworkInterfaces[0].MacAddress)
}

// NOTE: We expect there to be two tags on the instance to determine
// if it was in the latest scan or not. We need the interval of the scan
// such as 1d or 1w, etc. and then we expect the lastScanned tag to be in
// RFC822 format.
// TODO: Maybe the user would want to pass in their own format.
func (this *AwsInstance) IsInLatestScan() string {
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
	if lastScanned.After(Now().Add(interval * -1)) {
		return "Yes"
	} else {
		return "No"
	}
}

func (this *AwsInstance) LoadTags() {
	this.tags = make(map[string]string)
	for _, tag := range this.rawInstance.Tags {
		this.tags[aws.StringValue(tag.Key)] = aws.StringValue(tag.Value)
	}
}
