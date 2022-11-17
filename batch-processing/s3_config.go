package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func s3ConfigAndUpload(outputFile *os.File, row ProcessImage, outputErrRecords [][]string) (string, error) {
	awsRoleArn := os.Getenv("AWS_ROLE_ARN")
	if awsRoleArn == "" {
		log.Fatalln("Please set AWS_ROLE_ARN environment variable")
	}
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		log.Fatalln("Please set AWS_REGION environment variable")
	}
	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Fatalln("Please set S3_BUCKET environment variable")
	}

	config := &AWSConfig{
		ArnRole:    awsRoleArn,
		Region:     awsRegion,
		BucketName: s3Bucket,
	}

	// Set up S3 session
	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials.
	sess := session.Must(session.NewSession())

	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the ARN.
	creds := stscreds.NewCredentials(sess, config.ArnRole)

	// Create service client value configured for credentials
	// from assumed role.
	svc := s3.New(sess, &aws.Config{Credentials: creds})
	outputKey := filepath.Base(row.output)
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(outputKey),
		Body:   outputFile,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
			outputErrRecords = append(outputErrRecords, []string{row.url, row.input, row.output, err.Error()})
			fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
			outputErrRecords = append(outputErrRecords, []string{row.url, row.input, row.output, err.Error()})
		}
		return "", err
	}
	s3url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.BucketName, config.Region, outputKey)
	return s3url, nil
}
