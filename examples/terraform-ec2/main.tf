provider "aws" {
  region = "${var.region}"
}

resource "aws_s3_bucket" "fedrampup" {
  bucket = "${var.s3_bucket}"
  acl    = "private"

  versioning {
    enabled = true
  }
}

resource "aws_instance" "fedrampup" {
  ami           = "${data.aws_ami.amznlinux.id}"
  instance_type = "t2.small"
  user_data     = "${data.template_file.user_data.rendered}"
  key_name      = "deployer"
  security_groups = ["${aws_security_group.linux.id}"]
  iam_instance_profile = "${aws_iam_instance_profile.fedrampup.name}"
  subnet_id = "${var.subnet_id}"
  tags {
    Name = "fedrampup"
  }
}
