# FedRAMPup - Terraform EC2

Since at the time of this writing, ECS is not available in GovCloud, we'll make some terraform to do the following:

- Create a small EC2 instance for the FedRAMPup service that will run on a cron
- The output will be pushed to a versioned s3 bucket
- Set up IAM and Security Group rules as well
- Allow user to pass in VPC and Subnet Ids via TF_VAR variables
