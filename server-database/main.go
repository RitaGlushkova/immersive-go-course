package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

)
var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
var wordBarrierRegex = regexp.MustCompile(`(\w)([A-Z])`)

type Image struct {
	Title   string
	AltText string
	URL     string
}
type conventionalMarshaller struct {
	Value interface{}
}
func (c conventionalMarshaller) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(c.Value)

	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			return bytes.ToLower(wordBarrierRegex.ReplaceAll(
				match,
				[]byte(`${1}_${2}`),
			))
		},
	)

	return converted, err
}

func handlerImages(w http.ResponseWriter, r *http.Request) {
	images := []Image{
		{
			Title:   "Sunset",
			AltText: "Clouds at sunset",
			URL:     "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
		},
		{
			Title:   "Mountain",
			AltText: "A mountain at sunset",
			URL:     "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
		},
	}
	encoded, err := json.MarshalIndent(conventionalMarshaller{images}, "", "  ")
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}
	w.Write([]byte(encoded))
}

func main() {

	http.HandleFunc("/images.json", handlerImages)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
