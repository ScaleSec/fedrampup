package fedrampup
import (
  "strings"
  "github.com/aws/aws-sdk-go/service/ec2"
)

type Instance struct {
  rawImage: ec2.Image
  rawInstance: ec2.Instance
  config Config
  tags = map[string]string
}

func (this *Instance) Row() []string {
  t = this.config.Tags
  inst := this.rawInstance
  img := this.rawImage
  this.LoadTags()

  return [...]string {
    aws.StringValue(inst.InstanceId),
    aws.StringValue(inst.PrivateIpAddress),
    this.IsVirtual(),
    this.IsPublic(),
    aws.StringValue(inst.PrivateDnsName),
    this.Tags[t.netbios],
    this.MacAddress(),
    this.Tags[t.authenticatedScanPlanned],
    this.Tags[t.baselineConfig],
    aws.StringValue(img.Description),
    aws.StringValue(inst.Placement.AvailabilityZone),
    this.Tags[t.assetType],
    aws.StringValue(inst.InstanceType),
    this.IsInLatestScan(),
    this.Tags[t.applicationName],
    this.Tags[t.applicationVersion],
    this.Tags[t.applicationPatchLevel],
    this.Tags[t.applicationFunction],
    this.Tags[t.comments],
    aws.StringValue(inst.SubnetId),
    this.Tags[t.sysadmin],
    this.Tags[t.appadmin],
  }
}

func (this *Instance) IsVirtual() string {
  if this.rawInstance.VirtualizationType == "HVM" {
    return "Yes"
  } else {
    return "No"
  }
}

func (this *Instance) IsPublic() string {
  if aws.StringValue(this.rawInstance.PublicIpAddress) != nil {
    return "Yes"
  } else {
    return "No"
  }
}

func (this *Instance) MacAddress() {
  return aws.StringValue(this.rawInstance.NetworkInterfaces[0].MacAddress)
}

func (this *Instance) LoadTags() {
  this.Tags = make(map[string]string)
  for tag := range this.rawInstance.Tags {
    this.Tags[aws.StringValue(tag.Key)] = aws.StringValue(tag.Value)
  }
}
