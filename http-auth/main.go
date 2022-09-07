package main

import (
	"net/http"
	 "fmt"
	 "strings"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
	})
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
