package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/gographics/imagick.v2/imagick"
)

type Row struct {
	url    string
	input  string
	output string
}

type ProcessDownloadImage struct {
	url   string
	input string
	err   error
}
type ProcessUploadImage struct {
	url    string
	input  string
	output string
	err    error
}

type RowError struct {
	url string
	err error
}

type ConvertImageCommand func(args []string) (*imagick.ImageCommandResult, error)

type Converter struct {
	cmd ConvertImageCommand
}

func genFilepath(out, suffix string) string {
	if out != "" {
		return fmt.Sprintf("%d-%d-%s.%s", time.Now().UnixMilli(), rand.Int(), out, suffix)
	}
	return fmt.Sprintf("%d-%d.%s", time.Now().UnixMilli(), rand.Int(), suffix)
}

func (c *Converter) Grayscale(inputFilepath string, outputFilepath string) error {
	// Convert the image to grayscale using imagemagick
	// We are directly calling the convert command
	_, err := c.cmd([]string{
		"convert", inputFilepath, "-set", "colorspace", "Gray", outputFilepath,
	})
	return err
}

func ReadCsvFile(filename, headerTitle string) ([]string, error) {
	// Open CSV file to read url links
	fileContent, err := os.Open(filename)
	if err != nil {
		return []string{}, err
	}
	defer fileContent.Close()

	//read from this file
	records, err := csv.NewReader(fileContent).ReadAll()
	if err != nil {
		return []string{}, err
	}

	//check if there is any content
	if len(records) == 0 {
		return []string{}, fmt.Errorf("empty csv")
	}
	header := records[0]
	// find index of needed header if none return error
	//TODO processing headers

	//check if header and title info is what we expect
	// assuming we have multiple columns in the file, we need to find url
	var urls []string
	var counter int
	for i, h := range header {
		if h == headerTitle {
			fmt.Println(records)
			for n := 1; n < len(records); n++ {
				urls = append(urls, records[n][i])
				fmt.Println("read a record")
			}
		} else {
			counter++
		}
	}
	// no matching headers found
	// CHANGE
	if counter != len(header)-1 {
		return []string{}, fmt.Errorf("header not found, expected %s, got %v", headerTitle, header)
	}
	return urls, nil
}

func DownloadImages(urlsChan chan string, inputPathsChan chan ProcessDownloadImage, inputPath string, processingErrorChan chan RowError, wg *sync.WaitGroup) {
	for url := range urlsChan {
		d := DownloadImage(url, inputPath)
		if d.err != nil {
			processingErrorChan <- RowError{d.url, d.err}
			wg.Done()
		} else {
			inputPathsChan <- d
		}
	}
}

func ConvertImages(inputPathsChan chan ProcessDownloadImage, outputPath string, outputPathsChan chan Row, processingErrorChan chan RowError, wg *sync.WaitGroup) {
	for inputPath := range inputPathsChan {
		conv := ConvertImageIntoGreyScale(inputPath.input, outputPath, inputPath.url)
		if conv.err != nil {
			processingErrorChan <- RowError{url: conv.url, err: conv.err}
			wg.Done()
		} else {
			row := Row{
				url:    conv.url,
				input:  conv.input,
				output: conv.output,
			}
			outputPathsChan <- row
			wg.Done()
		}
	}
}
func DownloadImage(url, inputPath string) ProcessDownloadImage {
	start := time.Now()
	defer func() {
		fmt.Printf("downloaded file in %s\n", time.Since(start))
	}()
	//make GET request to URL
	r, err := http.Get(url)
	if err != nil {
		return ProcessDownloadImage{url: url, input: "", err: fmt.Errorf("couldn't process get request from %s. Error: %v", url, err)}
	}
	defer r.Body.Close()

	//create file where to download content of url
	inputFilepath := filepath.Join("/tmp", genFilepath("", "jpg"))

	// need to change to temp directory
	file, err := os.Create(inputFilepath)
	if err != nil {
		return ProcessDownloadImage{url: url, input: "", err: fmt.Errorf("file could not be created: %v", err)}
	}
	defer file.Close()

	//copy content from URl to the file
	_, err = io.Copy(file, r.Body)
	if err != nil {
		return ProcessDownloadImage{url: url, input: "", err: fmt.Errorf("data not copied into a file: %v", err)}
	}

	return ProcessDownloadImage{url: url, input: inputFilepath, err: nil}
}

func ConvertImageIntoGreyScale(inputFilepath, outputPath string, url string) ProcessUploadImage {
	// Set up imagemagick
	imagick.Initialize()
	defer imagick.Terminate()
	outputFilepath := filepath.Join("/tmp", genFilepath("out", "jpeg"))
	// Log what we're going to do
	log.Printf("processing: %q to %q\n", inputFilepath, outputFilepath)

	// Build a Converter struct that will use imagick
	c := &Converter{
		cmd: imagick.ConvertImageCommand,
	}

	// Do the conversion! and save to output folder
	err := c.Grayscale(inputFilepath, outputFilepath)
	if err != nil {
		return ProcessUploadImage{url: url, input: inputFilepath, output: "", err: fmt.Errorf("error: %v\n", err)}
	}
	// Log what we did
	log.Printf("processed: %q to %q\n", inputFilepath, outputFilepath)
	return ProcessUploadImage{url: url, input: inputFilepath, output: outputFilepath, err: nil}
}

func CreateAndWriteToCSVFile(path string, records [][]string) error {
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
