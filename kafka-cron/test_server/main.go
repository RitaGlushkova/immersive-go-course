package main

import (
	"io"
	"log"
	"net/http"
)

func getResponse(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hey I am working!\n")
}

func main() {
	http.HandleFunc("/", getResponse)
	log.Print("Starting server on :9009")
	err := http.ListenAndServe("localhost:9009", nil)
	log.Fatal(err)
}
