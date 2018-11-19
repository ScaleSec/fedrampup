package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"strings"
)

func main() {
	config := NewConfig()
	sess := GetSession()

	instances, err := NewFetcher(config, sess).Run()
	if err != nil {
		log.Fatal(err)
	}

	buffer := new(bytes.Buffer)
	w := csv.NewWriter(buffer)
	w.Write(Headers)
	for _, instance := range instances {
		if err := w.Write(instance.Row()); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()

	if strings.HasPrefix(config.OutputFile, "s3://") {
		path := strings.Split(config.OutputFile, "/")
		uploader := s3manager.NewUploader(sess)
		if _, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(path[2]),
			Key:    aws.String(strings.Join(path[3:], "/")),
			Body:   buffer,
		}); err != nil {
			log.Fatal(err)
		}
	} else {
		file, err := os.Create(config.OutputFile)
		if err != nil {
			fmt.Print(buffer.String())
			log.Fatal(err)
		}
		file.WriteString(buffer.String())
	}
}

func GetSession() *session.Session {
	var creds *credentials.Credentials
	sess := session.Must(session.NewSession())

	if len(os.Getenv("AWS_ACCESS_KEY_ID")) > 0 {
		creds = credentials.NewEnvCredentials()
	} else {
		creds = credentials.NewCredentials(&ec2rolecreds.EC2RoleProvider{})
	}
	if _, err := creds.Get(); err != nil {
		log.Fatal(err)
	}
	sess.Config = defaults.Config().WithCredentials(creds)
	return sess
}
