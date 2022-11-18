package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ssgreg/bottleneck"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type AWSConfig struct {
	ArnRole    string
	Region     string
	BucketName string
}

var bc *bottleneck.Calculator

func main() {

	// Accept --input and --output arguments for the images
	inputFilePath := flag.String("input", "", "A path to file with images to be processed")
	outputFilePath := flag.String("output", "", "A path to output file")
	outputPathFailed := flag.String("output-failed", "", "A path to file where filed outputs recorded")
	flag.Parse()
	// Ensure that all flags were set
	if *inputFilePath == "" || *outputFilePath == "" || *outputPathFailed == "" {
		flag.Usage()
		os.Exit(1)
	}
	bc = bottleneck.NewCalculator()
	imagick.Initialize()
	defer imagick.Terminate()
	urls, err := ReadCsvFile(*inputFilePath, "url")
	if err != nil {
		log.Fatal(err)
	}

	outputRecords := make([][]string, 0)
	outputErrRecords := make([][]string, 0)

	//append headers to slices, which we will use to write in our CSV files
	outputRecords = append(outputRecords, []string{"url", "input", "output", "s3url", "error"})
	outputErrRecords = append(outputErrRecords, []string{"url", "input", "output", "error"})

	//setting up WaitGroup to keep track of completion of our goroutines.
	var wg sync.WaitGroup

	// create channels
	// we can control how we process images.
	channels := Channels{
		urlsChan:            make(chan string, len(urls)),
		processingErrorChan: make(chan ProcessImage, len(urls)),
		inputPathsChan:      make(chan ProcessImage, len(urls)),
		outputPathsChan:     make(chan ProcessImage, len(urls))}

	// set go routines
	for i := 0; i < 4; i++ {
		go DownloadImages(channels, &wg)
		go ConvertImages(channels, &wg)
	}

	for _, url := range urls {
		wg.Add(1)
		channels.urlsChan <- url
	}

	for range urls {
		select {
		case invalidRecord := <-channels.processingErrorChan:
			fmt.Println(invalidRecord.suffix)
			outputErrRecords = append(outputErrRecords, []string{invalidRecord.url, invalidRecord.input, invalidRecord.output, invalidRecord.err.Error()})
		case row := <-channels.outputPathsChan:
			bc.TimeSlice(bottleneck.Index3)
			outputFile, err := os.Open(row.output)
			defer outputFile.Close()
			if err != nil {
				fmt.Printf("can not open file %v, err: %v\n", row.output, err)
				outputErrRecords = append(outputErrRecords, []string{row.url, row.input, row.output, err.Error()})
				continue
			}
			s3url, err := s3ConfigAndUpload(outputFile, row, outputErrRecords)
			if err != nil {
				continue
			}
			outputRecords = append(outputRecords, []string{row.url, row.input, row.output, s3url, "nil"})
		}
	}
	entries := bc.Stats()
	wg.Wait()
	fmt.Printf("Read file is %v, downloading image %v, converting %v\n, saving %v\n", entries[0].Duration, entries[1].Duration, entries[2].Duration, entries[3].Duration)
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
