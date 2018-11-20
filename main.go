package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
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

type Rowable interface {
	Row() []string
}

func main() {
	config := NewConfig()
	sess := GetSession()

	instances, err := NewFetcher(config, sess).Run()
	if err != nil {
		log.Fatal(err)
	}
	// To pass in the instances as a slice of interfaces you need
	// to manually convert to the interface.
	rowables := make([]Rowable, len(instances))
	for idx, inst := range instances {
		rowables[idx] = inst
	}
	buffer := Dump(rowables, config)
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
		defer file.Close()
		if err != nil {
			fmt.Print(buffer.String())
			log.Fatal(err)
		}
		file.WriteString(buffer.String())
	}
}

func Dump(instances []Rowable, config *Config) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	if config.OutputFormat() == "csv" {
		w := csv.NewWriter(buffer)
		w.Write(Headers)
		for _, instance := range instances {
			if err := w.Write(instance.Row()); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}
		w.Flush()
	} else if config.OutputFormat() == "json" {
		jsonData, err := json.MarshalIndent(instances, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		if _, err := buffer.Write(jsonData); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("Unsupported output format %s", config.OutputFormat())
	}
	return buffer
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
