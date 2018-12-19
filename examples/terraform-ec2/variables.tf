variable "s3_bucket" {}
variable "go_version" { default = "go1.11.2" }
variable "s3_path" { default = "output.csv" }
variable "region" { default = "us-gov-west-1" }
variable "vpc_id" {}
variable "subnet_id" {}
locals {
  "s3_uri" = "s3://${var.s3_bucket}/${var.s3_path}"
}
