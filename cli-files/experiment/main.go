package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	//fmt.Println("Hello")
	file, err := os.Open("rita")
	if err != nil {
		log.Fatal(err)
	}
	buffer := make([]byte, 8)
	var result []byte
// how to use one slice rather than 2
//re-allocation
	for {
		n, err := file.Read(buffer)
		result = append(result, buffer[0:n]...)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(string(result))
}
