package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/gographics/imagick.v2/imagick"
)

type Row struct {
	url string
	title string
}

type ConvertImageCommand func(args []string) (*imagick.ImageCommandResult, error)

type Converter struct {
	cmd ConvertImageCommand
}

func (c *Converter) Grayscale(inputFilepath string, outputFilepath string) error {
	// Convert the image to grayscale using imagemagick
	// We are directly calling the convert command
	_, err := c.cmd([]string{
		"convert", inputFilepath, "-set", "colorspace", "Gray", outputFilepath,
	})
	return err
}

func ReadCsvFile(filename, headerTitle string) ([][]string, error) {
   // Open CSV file from there to read url links
   fileContent, err := os.Open(filename)
   if err != nil {
      return [][]string{}, err
   }
   defer fileContent.Close()

   //read from this file
   records, err := csv.NewReader(fileContent).ReadAll()
   if err != nil {
      return [][]string{}, err
   }

   //check if there is any content
   if len(records) == 0 {
	return [][]string{}, fmt.Errorf("empty csv")
   }

   //check if header and thus info is what we expect
   header := records[0]
   if header[0] != headerTitle || len(header) == 0 {
	return [][]string{}, fmt.Errorf("incorrect header, expected %s, got %s", headerTitle, header[0])
   }

   //if no errors return records (first will be title)
   return records, nil
}

func (row *Row) DownloadAndSaveImage(url, name, inputPath, outputPath string) error{

	//make GET request to URL
    r, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("couldn't process get request from %s. Error: %v", url, err)
    }
    defer r.Body.Close()

	//create file where to download content of url
	fileName := fmt.Sprintf("%s.jpg", name)
	inputFilepath := filepath.Join(inputPath, fileName)
	outputFilepath := filepath.Join(outputPath, fileName)
	file, err := os.Create(inputFilepath)
    if err != nil {
        return err
    }
    defer file.Close()

	//copy content from URl to the file
    _, err = io.Copy(file, r.Body)
    if err != nil {
        return err
	}

	// Set up imagemagick
	imagick.Initialize()
	defer imagick.Terminate()

	 // Log what we're going to do
	 log.Printf("processing: %q to %q\n", inputFilepath, outputFilepath)

	// Build a Converter struct that will use imagick
	c := &Converter{
		cmd: imagick.ConvertImageCommand,
	}

	// Do the conversion! and save to output folder
	err = c.Grayscale(inputFilepath, outputFilepath)
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	 // Log what we did
 	log.Printf("processed: %q to %q\n", inputFilepath, outputFilepath)

	return nil		

}

