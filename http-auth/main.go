package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

var htmlHead = "<!DOCTYPE html><html>"

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		w.Header().Set("Content-Type", "text/html") //text/plain sends a string as response
		keys, ok := r.URL.Query()["foo"]
		if ok {
			foo := keys[0]
			w.Write([]byte(fmt.Sprintf("%v<em>Hello, world</em><p>Query parameters:<ul><li>foo: %v</li></ul>\n", htmlHead, html.EscapeString(foo))))
		} else {
			w.Write([]byte(fmt.Sprintf("%v<em>Hello, world</em>\n", htmlHead)))
		}

	case "POST":
		w.Header().Set("Content-Type", "text/html")
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(fmt.Sprintf("%v%v\n", htmlHead, html.EscapeString(string(b)))))
	default:
		w.Write([]byte("Only GET and POST methods are available"))
	}
}
func handler200(w http.ResponseWriter, r *http.Request) {
	prefix := "/"
	w.Write([]byte(fmt.Sprintf("%v\n", strings.TrimPrefix(r.URL.Path, prefix))))
}
func handler500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	message := "Internal server error"
	w.Write([]byte(fmt.Sprintf("%v\n", message)))
}

func handlerAuth(w http.ResponseWriter, r *http.Request) {
	username := goDotEnvVariable("AUTH_USERNAME")
	password := goDotEnvVariable("AUTH_PASSWORD")
	u, p, ok := r.BasicAuth()
	if !ok || u != username || p != password {
		fmt.Println("Error parsing basic auth")
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		w.WriteHeader(401)
		return
	}
	w.Write([]byte(fmt.Sprintf("%v<em>Hello, %s!</em>\n", htmlHead, username)))
}
func handlerLimit(w http.ResponseWriter, r *http.Request) {
	limiter := rate.NewLimiter(100, 30)
	if limiter.Allow() {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf("%v<em>Hello, world</em>\n", htmlHead)))
	} else {
		w.WriteHeader(503)
		message := "503 Service Unavailable"
		w.Write([]byte(fmt.Sprintf("%v\n", message)))
	}
}

func main() {
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/200", handler200)
	http.HandleFunc("/authenticated", handlerAuth)
	http.HandleFunc("/500", handler500)
	http.HandleFunc("/404", http.NotFoundHandler().ServeHTTP)
	http.HandleFunc("/limited", handlerLimit)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
