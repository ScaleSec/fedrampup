data "aws_iam_policy_document" "fedrampup" {
  statement {
    actions   = ["s3:PutObject"]
    resources = [
      "${aws_s3_bucket.fedrampup.arn}",
      "${aws_s3_bucket.fedrampup.arn}/*",
    ]
  }
}

resource "aws_iam_policy" "fedrampup" {
  name = "FedrampupBucketWrite"
  policy = "${data.aws_iam_policy_document.fedrampup.json}"
}

resource "aws_iam_role" "fedrampup" {
  name = "fedrampup"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "fedrampup_attachment" {
  role       = "${aws_iam_role.fedrampup.name}"
  policy_arn = "${aws_iam_policy.fedrampup.arn}"
}

resource "aws_iam_instance_profile" "fedrampup" {
  name = "fedrampup_profile"
  role = "${aws_iam_role.fedrampup.name}"
}
