package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSConfig struct {
	ArnRole    string
	Region     string
	BucketName string
}

func main() {

	// Accept --input and --output arguments for the images
	inputFilePath := flag.String("input", "", "A path to an image to be processed")
	outputFilePath := flag.String("output", "", "A path to where the processed image should be written")
	outputPathFailed := flag.String("output-failed", "", "A path to where failed image should be written")
	flag.Parse()

	// Ensure that both flags were set
	if *inputFilePath == "" || *outputFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}
	records, err := ReadCsvFile(*inputFilePath, "url")
	if err != nil {
		log.Fatal(err)
	}
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

	outputRecords := make([][]string, 0)
	outputErrRecords := make([][]string, 0)
	outputRecords = append(outputRecords, []string{"url", "input", "output", "s3url"})
	outputErrRecords = append(outputErrRecords, []string{"url", "err", "message"})
	records = records[1:]
	var wg sync.WaitGroup
	urlsChan := make(chan string, 4)
	processingErrorChan := make(chan RowError)
	inputPathsChan := make(chan ProcessDownloadImage, 4)
	outputPathsChan := make(chan Row, 4)
	go DownloadImageS(urlsChan, inputPathsChan, *inputFilePath, processingErrorChan, &wg)
	go ConvertImages(inputPathsChan, *outputFilePath, outputPathsChan, processingErrorChan, &wg)

	for _, record := range records {
		wg.Add(1)
		urlsChan <- record[0]
		select {
		case invalidRecord := <-processingErrorChan:
			fmt.Println(invalidRecord.message)
			outputErrRecords = append(outputErrRecords, []string{invalidRecord.url, invalidRecord.err.Error(), invalidRecord.message})
		case row := <-outputPathsChan:
			outputFile, err := os.Open(row.output)
			if err != nil {
				fmt.Printf("can not open file %v, err: %v\n", row.output, err)
				break
			}

			// Uploads the object to S3. The Context will interrupt the request if the
			// timeout expires.
			outputKey := filepath.Base(row.output)
			_, err = svc.PutObject(&s3.PutObjectInput{
				Bucket: aws.String(config.BucketName),
				Key:    aws.String(outputKey),
				Body:   outputFile,
			})
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
					// If the SDK can determine the request or retry delay was canceled
					// by a context the CanceledErrorCode error code will be returned.
					fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
				} else {
					fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
				}
				os.Exit(1)
			}
			s3url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.BucketName, config.Region, outputKey)
			outputRecords = append(outputRecords, []string{row.url, row.input, row.output, s3url})
		}
	}
	wg.Wait()

	err = CreateFile(*outputFilePath, outputRecords)
	if err != nil {
		log.Fatal(err)
	}

	err = CreateFile(*outputPathFailed, outputErrRecords)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateFile(path string, records [][]string) error {
	csvFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create a file %s, %v", path, err)
	}
	defer csvFile.Close()
	w := csv.NewWriter(csvFile)
	defer w.Flush()
	err = w.WriteAll(records)
	if err != nil {
		return fmt.Errorf("failed to write records into file %s, %v", path, err)
	}
	return nil
}
