package main

import (
	"encoding/csv"
	"fmt"
	"os"
)


func ReadCsvFile(filename string) ([][]string, error) {
   // Open CSV file
   fileContent, err := os.Open(filename)
   if err != nil {
      return [][]string{}, err
   }
   defer fileContent.Close()

   // Read File into a Variable
   records, err := csv.NewReader(fileContent).ReadAll()

   if err != nil {
      return [][]string{}, err
   }

   if len(records) == 0 {
	return [][]string{}, fmt.Errorf("empty csv")
   }

   header := records[0]

   if header[0] != "url" || len(header) == 0 {
	return [][]string{}, fmt.Errorf("incorrect header %s, url is required", header[0])
   }

   return records, nil
}

// func downloadAnImage(url string, name string) error{
//     r, err := http.Get(url)
//     if err != nil {
//         return fmt.Errorf("couldn't process get request from %s. Error: %v", url, err)
//     }
//     defer r.Body.Close()

//     fname := name+".jpg"
// 	//create file where to download the image
//     f, err := os.Create(fname)
//     if err != nil {
//         return err
//     }
//     defer f.Close()

//     _, err = io.Copy(fname, r.Body)

//     if err != nil {
//         return err
// }
// fmt.Println("image downloaded")
// 	return nil
// }