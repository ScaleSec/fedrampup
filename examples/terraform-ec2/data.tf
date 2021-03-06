data "aws_caller_identity" "current" {}

data "aws_ami" "amznlinux" {
  most_recent = true

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-2.0.20181114-x86_64-gp2"]
  }
}

data "template_file" "user_data" {
  template = "${file("user-data.tpl")}"

  vars {
    go_version = "${var.go_version}"
    aws_region = "${var.region}"
    s3_uri  = "${local.s3_uri}"
  }
}

data "aws_iam_policy" "security_audit" {
  arn = "arn:aws-us-gov:iam::aws:policy/SecurityAudit"
}
