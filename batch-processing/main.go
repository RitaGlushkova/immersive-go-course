package main

import (
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
	inputFilePath := flag.String("input", "", "A path to file with images to be processed")
	outputFilePath := flag.String("output", "", "A path to output file")
	outputPathFailed := flag.String("output-failed", "", "A path to file where filed outputs recorded")
	flag.Parse()

	// Ensure that all flags were set
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

	//append headers to slices, which we will use to write in our CSV files
	outputRecords = append(outputRecords, []string{"url", "input", "output", "s3url"})
	outputErrRecords = append(outputErrRecords, []string{"url", "err", "message"})

	//setting up WaitGroup to keep track of completion of our goroutines.
	var wg sync.WaitGroup

	// create channels
	// we can control how we process images.
	urlsChan := make(chan string, len(records))
	processingErrorChan := make(chan RowError, len(records))
	inputPathsChan := make(chan ProcessDownloadImage, len(records))
	outputPathsChan := make(chan Row, len(records))

	// set go routines
	for i := 0; i < 4; i++ {
		go DownloadImages(urlsChan, inputPathsChan, *inputFilePath, processingErrorChan, &wg)
		go ConvertImages(inputPathsChan, *outputFilePath, outputPathsChan, processingErrorChan, &wg)
	}
	for _, record := range records {
		wg.Add(1)
		urlsChan <- record //url
	}

	for range records {
		select {
		case invalidRecord := <-processingErrorChan:
			fmt.Println(invalidRecord.message)
			outputErrRecords = append(outputErrRecords, []string{invalidRecord.url, invalidRecord.err.Error(), invalidRecord.message})
		case row := <-outputPathsChan:
			outputFile, err := os.Open(row.output)
			defer outputFile.Close()
			if err != nil {
				fmt.Printf("can not open file %v, err: %v\n", row.output, err)
				break
			}
			outputKey := filepath.Base(row.output)
			_, err = svc.PutObject(&s3.PutObjectInput{
				Bucket: aws.String(config.BucketName),
				Key:    aws.String(outputKey),
				Body:   outputFile,
			})
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
					fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
					outputErrRecords = append(outputErrRecords, []string{row.url, err.Error(), "upload canceled due to timeout"})
				} else {
					fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
					outputErrRecords = append(outputErrRecords, []string{row.url, err.Error(), "failed to upload object to S3 bucket"})
				}
				continue
			}
			s3url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.BucketName, config.Region, outputKey)
			outputRecords = append(outputRecords, []string{row.url, row.input, row.output, s3url})
		}
	}
	wg.Wait()

	errOutput := CreateAndWriteToCSVFile(*outputFilePath, outputRecords)
	if errOutput != nil {
		fmt.Fprintf(os.Stderr, "failed to creates and write to the output.csv file, path: %v, %v\n", *outputFilePath, err)
	}

	errFailed := CreateAndWriteToCSVFile(*outputPathFailed, outputErrRecords)
	if errFailed != nil {
		fmt.Fprintf(os.Stderr, "failed to creates and write to the failed.csv file, path: %v, %v\n", *outputPathFailed, err)
	}

	if errOutput != nil || errFailed != nil || len(outputErrRecords) > 1 {
		os.Exit(1)
	}
}
