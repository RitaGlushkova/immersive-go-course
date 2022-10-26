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
	//fmt.Println(len(records))
	if err != nil {
		log.Fatal(err)
	}
	outputRecords := make([][]string, len(records))
	outputRecords = append(outputRecords, []string{"url", "input", "output"})
	records = records[1:]
	var wg sync.WaitGroup
	urlsChan := make(chan string)
	processingErrorChan := make(chan error)
	inputPathsChan := make(chan ProcessDownloadImage)
	outputPathsChan := make(chan ProcessUploadImage)
	//for i:=0; i<3; i++ {
	go DownloadImageS(urlsChan, inputPathsChan, *inputPath, processingErrorChan, &wg)
	//}
	go ConvertImages(inputPathsChan, *outputPath, outputPathsChan, processingErrorChan, &wg)

	for _, record := range records {
		wg.Add(1)
		urlsChan <- record[0]
	}
	wg.Wait()

	select {
	case err := <-processingErrorChan:
		fmt.Println(err)
	case row := <-outputPathsChan:
		outputRecords = append(outputRecords, []string{row.url, row.input, row.output})
	}

	csvFile, err := os.Create(filepath.Join(*outputPath, "output.csv"))
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFile.Close()
	w := csv.NewWriter(csvFile)
	err = w.WriteAll(outputRecords)
	if err != nil {
		log.Fatal(err)
	}
}
