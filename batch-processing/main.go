package main

import (
	"encoding/csv"
	"flag"
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

	records, err := ReadCsvFile("/inputs/input.csv", "url")
	if err != nil {
		log.Fatal(err)
	}
	outputRecords := make([][]string, len(records)-1)
	outputRecords = append(outputRecords, []string{"url", "input", "output"})
	records = records[1:]
	var wg sync.WaitGroup
	urlsChan := make(chan string, 3)
	inputPathsChan := make(chan ProcessDownloadImage, 3)
	outputPathsChan := make(chan ProcessUploadImage)
	outputWriteChan := make(chan Row)
	for i:=0; i<3; i++ {
	go DownloadImageS(urlsChan, inputPathsChan, *inputPath)
	}
	go ConvertImages(inputPathsChan, *outputPath, outputPathsChan)
	go WriteIntoOutputSlice(outputPathsChan, outputWriteChan, &wg)

	for _, record := range records {
		wg.Add(1)
		urlsChan <- record[0]
    }
    
	close(urlsChan)
	wg.Wait()

	for row := range outputWriteChan {
		outputRecords = append(outputRecords, []string{row.url, row.input, row.output})
    }

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
}
