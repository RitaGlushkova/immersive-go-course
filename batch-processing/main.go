package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	outputRecords = append(outputRecords, []string{"url", "title"})
	records = records[1:]
	for i, record := range records {

		row := Row {
			url: record[0],
			title: fmt.Sprintf("%s.%s", fmt.Sprint(i), "jpg"),
		}
		err := row.DownloadAndSaveImage(record[0], fmt.Sprint(i), *inputPath, *outputPath)
		if err != nil {
			log.Fatal(err)
		}
		outputRecords = append(outputRecords, []string{row.url, row.title})
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
