#!/bin/bash
# Dependencies
yum install -y git

# Download and install golang
VERSION=${go_version}
curl -o /tmp/$VERSION.tar.gz "https://dl.google.com/go/$VERSION.linux-amd64.tar.gz"
tar -C /usr/local -xzf /tmp/$VERSION.tar.gz

# Setup Go path
export GOPATH=/opt/go
mkdir -p /opt/go/src /opt/go/pkg /opt/go/bin

# Get the FedRAMPup package
/usr/local/go/bin/go get github.com/ScaleSec/fedrampup

# Create cron script to use ENV vars
WRAPPER=/opt/fedrampup-wrapper
cat << EOF > $WRAPPER
#!/bin/bash
export AWS_REGION=${aws_region}
export OUTPUT_FILE=${s3_uri}

/opt/go/bin/fedrampup
EOF

chmod +x $WRAPPER

# Add crontab
(crontab -l ; echo "00 00 * * * $WRAPPER") | crontab
