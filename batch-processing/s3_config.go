package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func s3Config() (svc *s3.S3, config *AWSConfig) {
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

	config = &AWSConfig{
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
	svc = s3.New(sess, &aws.Config{Credentials: creds})

	return svc, config
}
