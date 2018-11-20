# FedRAMPup!

How annoying is it to get an inventory of AWS resources for FedRAMP and then put that into the SSP formatted CSV? Very. `fedrampup` will handle all of this for you. It can be run on it's own from a Docker container or in AWS in Fargate with S3 output.

# Run

With Go installed

```
go get -u https://github.com/ScaleSec/fedrampup

# Assuming your environment variables are configured correctly just run
fedrampup
```


Without Go installed, but using Docker

```
# Example:
docker run -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
           -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
           -e "OUTPUT_FILE=s3://my-bucket/output.csv" fedrampup
```

# Authentication

If you pass in [AWS environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html) into `fedrampup` it will use those, otherwise it will assume you are running on an EC2 instance or ECS container (recommended) and use standard instance profile authentication.

# Configuration

Everyone's environment is different. That's why `fedrampup` has an extensive configuration framework based in environments in keeping with a [12 factor application](https://12factor.net/).

FedRAMP SSP expects the assets in a CSV with the rows:

- Unique Asset Identifier
- IPv4 or IPv6 Address
- Virtual
- Public
- DNS Name or URL
- NetBIOS Name
- MAC Address
- Authenticated Scan
- Baseline Configuration Name
- OS Name and Version
- Location
- Asset Type
- Hardware Make/Model
- In Latest Scan
- Software/Database Name
- Software/Database Version
- Patch Level
- Function
- Comments
- Serial Number/Asset Tag Number
- VLAN/Network ID
- System Administrator/Owner
- Application Administrator/Owner

There are several things you can configure, mainly what tags you use to identify assets. A full table of environment variables that can be passed in is below:

| Env Var | Default Value | Description |
| --- | --- | --- |
|REGIONS | `us-gov-west-1` | Comma separated list of AWS regions |
|ROLES | | Comma separated list of AWS Role ARNs of separate accounts that will be scanned. If empty only the account with the EC2 Role for the host this is running on will be run. If not running on EC2, it will take credentials from Env |
|OUTPUT_FILE | `output.csv` | Output file name. If this starts with `s3://` it will be treated as an S3 URI. Output format can be one of: csv, json |
|AWS_REGION | `us-gov-west-1` | If S3 output is desired, this is the region of the S3 bucket.  
|SCAN_INTERVAL | `24h` | How often security scans are run in your organization in [Go duration format](https://golang.org/pkg/time/#ParseDuration) (i.e. 1h10m, 5m, etc.) |
|TAG_NETBIOS | `NetBIOS`| EC2 tag used for NetBIOS name for Windows hosts|
|TAG_ASSET_TYPE | `AssetType` | EC2 tag used for Asset Type |
|TAG_BASELINE_CONFIG | `BaselineConfiguration` | EC2 tag used for Baseline Configuration|
|TAG_AUTHENTICATED_SCAN_PLANNED | `AuthenticatedScanPlanned` | EC2 tag used for if and authenticated scan is planned |
|TAG_LAST_SCANNED | `LastScanned` | EC2 tag used for when the last scan was run on the host. Should be in RFC822 Format: "02 Jan 06 15:04 MST" (https://www.w3.org/Protocols/rfc822/ "Date and Time Specification") |
|TAG_APPLICATION_VENDOR | `ApplicationVendor` | EC2 tag used for Software/Database Vendor Name|
|TAG_APPLICATION_NAME_VERSION | `ApplicationNameAndVersion` | EC2 tag used for Software/Database Application Name and Version|
|TAG_APPLICATION_PATCH_LEVEL | `ApplicationPatchLevel` | EC2 tag used for Software/Database Patch Level |
|TAG_APPLICATION_FUNCTION | `ApplicationFunction` | EC2 tag used for Software/Database Function|
|TAG_COMMENTS | `Comments` | EC2 tag used for Comments|
|TAG_SERIAL_NUMBER | `SerialNumber` | EC2 tag used for the asset's serial number|
|TAG_SYSADMIN | `SysAdmin` | EC2 tag used for the asset's sysadmin or team|
|TAG_APPADMIN | `AppAdmin` | EC2 tag used for the asset's appadmin or team|
