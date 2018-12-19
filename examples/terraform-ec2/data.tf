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
