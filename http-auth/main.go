package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET": 
			w.Header().Set("Content-Type", "text/html") //text/plain sends a string as response
			w.Write([]byte("<!DOCTYPE html><html><em>Hello, world</em>\n"))
		case "POST": 
		w.Header().Set("Content-Type", "text/html")
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(b))
		default:
		w.Write([]byte("Only GET and POST methods are available"))
	}})
	
	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		prefix := "/"
		w.Write([]byte(fmt.Sprintf("%v\n", strings.TrimPrefix(r.URL.Path, prefix))))
	})
	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		message := "Internal server error"
		w.Write([]byte(fmt.Sprintf("%v\n", message)))
	})
	http.HandleFunc("/404", http.NotFoundHandler().ServeHTTP)
	http.ListenAndServe(":8080", nil)
}
