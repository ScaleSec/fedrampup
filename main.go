package main

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Rowable interface {
	Row() []string
}

func main() {
	config := NewConfig()
	sess := GetSession()
	journaler, err := NewJournaler()
	HandleError(err)

	defer journaler.Close()
	journaler.Write(Headers)

	filename, err := NewFetcher(config, sess, journaler).Run()
	HandleError(err)

	if strings.HasPrefix(config.OutputFile, "s3://") {
		buffer, err := ioutil.ReadFile(filename)
		HandleError(err)

		path := strings.Split(config.OutputFile, "/")
		uploader := s3manager.NewUploader(sess)
		// TODO: Better way to upload file from disk as opposed to memory
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(path[2]),
			Key:    aws.String(strings.Join(path[3:], "/")),
			Body:   bytes.NewBuffer(buffer),
		})
		HandleError(err)
	} else {
		// Since we've been journalling in CSV format we can just Rename
		// the journal to the output filename
		os.Rename(filename, config.OutputFile)
	}
}

func GetSession() *session.Session {
	var creds *credentials.Credentials
	sess := session.Must(session.NewSession())
	meta := ec2metadata.New(sess)
	if meta.Available() {
		creds = credentials.NewCredentials(&ec2rolecreds.EC2RoleProvider{
			Client: meta,
		})
	} else {
		creds = credentials.NewEnvCredentials()
	}
	if _, err := creds.Get(); err != nil {
		log.Fatal(err)
	}
	sess.Config = defaults.Config().WithCredentials(creds)
	return sess
}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
