package main

// import (
// 	"flag"
// 	"fmt"
// 	"log"
// 	"os"
// 	"sync"
// )

// var batchSize = 1
// func main() {
	
// 	// Accept --input and --output arguments for the images
// 	inputPath := flag.String("input", "", "A path to an image to be processed")
// 	outputPath := flag.String("output", "", "A path to where the processed image should be written")
// 	flag.Parse()

// 	// Ensure that both flags were set
// 	if *inputPath == "" || *outputPath == "" {
// 		flag.Usage()
// 		os.Exit(1)
// 	}

// 	records, err := ReadCsvFile("/inputs/input.csv", "url")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	records = records[1:]
// 	var wg sync.WaitGroup
// 	urlChan := make(chan string)
// 	lenURL := make(chan int)
// 	go convertToGreyscale(lenURL, &wg)
// 	for i:=0; i<5; i++ {
// 		go downloadImages(urlChan, lenURL)
// 	}
// 	for _, record := range records {
// 		wg.Add(1)
// 		urlChan <- record[0]
// 	}
// 	close(urlChan)
// 	wg.Wait()
	
// }

// func downloadImages(urls chan string, lens chan int) {
// 	for url := range urls {
// 	len := downloadImage(url)
// 	lens <- len
// 	}
// }

// func downloadImage(url string) int{
// 	fmt.Println(url)
// 	return len(url)
// }

// func convertToGreyscale(lens chan int, wg *sync.WaitGroup) {
// 	for len := range lens{
// 		fmt.Println(len)
// 			wg.Done()
// 	}
// }

// func addToCsv() {
	
// }