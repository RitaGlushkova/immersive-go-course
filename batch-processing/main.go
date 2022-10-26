package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func main() {

	// Accept --input and --output arguments for the images
	inputPath := flag.String("input", "", "A path to an image to be processed")
	outputPath := flag.String("output", "", "A path to where the processed image should be written")
	flag.Parse()

	// Ensure that both flags were set
	if *inputPath == "" || *outputPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	records, err := ReadCsvFile(filepath.Join(*inputPath, "input.csv"), "url")
	if err != nil {
		log.Fatal(err)
	}
	outputRecords := make([][]string, 0)
	outputErrRecords := make([][]string, 0)
	outputRecords = append(outputRecords, []string{"url", "input", "output"})
	outputErrRecords = append(outputErrRecords, []string{"url", "err", "message"})
	records = records[1:]
	var wg sync.WaitGroup
	urlsChan := make(chan string, 4)
	processingErrorChan := make(chan RowError)
	inputPathsChan := make(chan ProcessDownloadImage, 4)
	outputPathsChan := make(chan Row, 4)
	go DownloadImageS(urlsChan, inputPathsChan, *inputPath, processingErrorChan, &wg)
	go ConvertImages(inputPathsChan, *outputPath, outputPathsChan, processingErrorChan, &wg)

	for _, record := range records {
		wg.Add(1)
		urlsChan <- record[0]
		select {
		case invalidRecord := <-processingErrorChan:
			fmt.Println(invalidRecord.message)
			outputErrRecords = append(outputErrRecords, []string{invalidRecord.url, invalidRecord.err.Error(), invalidRecord.message})
		case row := <-outputPathsChan:
			outputRecords = append(outputRecords, []string{row.url, row.input, row.output})
		}
	}
	wg.Wait()
	csvFile, err := os.Create(filepath.Join(*outputPath, "output.csv"))
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFile.Close()
	w := csv.NewWriter(csvFile)
	defer w.Flush()
	err = w.WriteAll(outputRecords)
	if err != nil {
		log.Fatal(err)
	}

	csvFileErr, err := os.Create(filepath.Join(*outputPath, "failed.csv"))
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFileErr.Close()
	w = csv.NewWriter(csvFileErr)
	defer w.Flush()
	err = w.WriteAll(outputErrRecords)
	if err != nil {
		log.Fatal(err)
	}
}
