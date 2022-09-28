package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Image struct {
	Title   string `json:"title"`
	AltText string `json:"alt_text"`
	URL     string `json:"url"`
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
	//b, err := json.Marshal(images)
	b, err := json.MarshalIndent(images, "", " ")
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}
	w.Write([]byte(b))
}

func main() {

	http.HandleFunc("/images.json", handlerImages)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
