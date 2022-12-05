package main

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/exp/slices"
	"gopkg.in/gographics/imagick.v2/imagick"
)

type ProcessImage struct {
	url    string
	input  string
	output string
	err    error
	suffix string
}

type Channels struct {
	urlsChan            chan string
	processingErrorChan chan ProcessImage
	inputPathsChan      chan ProcessImage
	outputPathsChan     chan ProcessImage
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
		return []string{}, fmt.Errorf("couldn't open input file, %v", err)
	}
	defer fileContent.Close()

	//read from this file
	records, err := csv.NewReader(fileContent).ReadAll()
	if err != nil {
		return []string{}, fmt.Errorf("couldn't read input file, %v", err)
	}

	//check if there is any content
	if len(records) == 0 {
		return []string{}, fmt.Errorf("empty csv")
	}
	header := records[0]
	// find index of needed header if none return error

	//check if header and title info is what we expect
	// assuming we have multiple columns in the file, we need to find url
	var urls []string
	if header[0] == headerTitle {
		fmt.Println(records)
		for n := 1; n < len(records); n++ {
			urls = append(urls, records[n][0])
			fmt.Println("read a record")
		}
	} else {
		return []string{}, fmt.Errorf("no matching header found, expected %s, got %v", headerTitle, header)
	}
	return urls, nil
}

func DownloadImages(channels Channels, wg *sync.WaitGroup) {
	for url := range channels.urlsChan {
		d := DownloadImage(url)
		if d.err != nil {
			channels.processingErrorChan <- d
			wg.Done()
		} else {
			channels.inputPathsChan <- d
		}
	}
}

func ConvertImages(channels Channels, wg *sync.WaitGroup) {
	for inputPath := range channels.inputPathsChan {
		conv := ConvertImageIntoGreyScale(inputPath)
		if conv.err != nil {
			channels.processingErrorChan <- conv
		} else {
			channels.outputPathsChan <- conv
		}
		wg.Done()
	}
}
func DownloadImage(url string) ProcessImage {
	start := time.Now()
	defer func() {
		fmt.Printf("downloaded file in %s\n", time.Since(start))
	}()
	var supportedFormats = []string{"jpeg", "png", "gif"}
	//make GET request to URL
	r, err := http.Get(url)
	if err != nil || r.StatusCode != 200 {
		return ProcessImage{url: url, input: "no filepath", output: "no filepath", err: fmt.Errorf("couldn't fetch image. Error: %v", err)}
	}
	defer r.Body.Close()

	var buf bytes.Buffer

	tee := io.TeeReader(r.Body, &buf)
	_, suffix, err := image.DecodeConfig(tee)

	if err != nil {
		return ProcessImage{url: url, input: "no filepath", output: "no filepath", err: fmt.Errorf("couldn't decode image. Error: %v", err)}
	}
	//suffix := "jpeg"

	if !slices.Contains(supportedFormats, suffix) {
		return ProcessImage{url: url, input: "no filepath", output: "no filepath", err: fmt.Errorf("format: %s is not supported. Error: %v", suffix, err)}
	}
	//create file where to download content of url
	inputImagePath := filepath.Join("/tmp", genFilepath("", suffix))
	// need to change to temp directory
	file, err := os.Create(inputImagePath)
	if err != nil {
		return ProcessImage{url: url, input: inputImagePath, output: "no filepath", err: fmt.Errorf("file could not be created: %v", err), suffix: suffix}
	}
	defer file.Close()

	//copy content from URl to the file
	_, err = io.Copy(file, &buf)
	if err != nil {
		return ProcessImage{url: url, input: inputImagePath, output: "", err: fmt.Errorf("data not copied into a file: %v", err)}
	}

	return ProcessImage{url: url, input: inputImagePath, output: "", err: nil, suffix: suffix}
}

func ConvertImageIntoGreyScale(in ProcessImage) ProcessImage {
	outputImagePath := filepath.Join("/tmp", genFilepath("out", in.suffix))
	// Log what we're going to do
	log.Printf("processing: %q to %q\n", in.input, outputImagePath)

	// Build a Converter struct that will use imagick
	c := &Converter{
		cmd: imagick.ConvertImageCommand,
	}

	// Do the conversion! and save to output folder
	err := c.Grayscale(in.input, outputImagePath)
	if err != nil {
		return ProcessImage{url: in.url, input: in.input, output: outputImagePath, err: fmt.Errorf("Couldn't turn into grey scale: %v\n", err), suffix: in.suffix}
	}
	// Log what we did
	log.Printf("processed: %q to %q\n", in.input, outputImagePath)
	return ProcessImage{url: in.url, input: in.input, output: outputImagePath, err: nil, suffix: in.suffix}
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
