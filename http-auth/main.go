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
		defaultResponse := htmlHead+"<em>Hello, world</em>"
		if ok {
			foo := keys[0]
			defaultResponse = defaultResponse + fmt.Sprintf("<p>Query parameters:<ul><li>foo: %v</li></ul>\n", html.EscapeString(foo))
		} 
		w.Write([]byte(fmt.Sprintf("%v\n", defaultResponse)))
		
	case "POST":
		w.Header().Set("Content-Type", "text/html")
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			w.WriteHeader(500)
		}
		w.WriteHeader(202)
		w.Write([]byte(fmt.Sprintf("%v%v\n", htmlHead, html.EscapeString(string(b)))))
	default:
		w.WriteHeader(405)
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
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		w.WriteHeader(401)
		return
	}
	w.Write([]byte(fmt.Sprintf("%v<em>Hello, %s!</em>\n", htmlHead, username)))
}

func handlerLimit(limiter *rate.Limiter, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			w.WriteHeader(429)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

func main() {
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/200", handler200)
	http.HandleFunc("/authenticated", handlerAuth)
	http.HandleFunc("/500", handler500)
	http.HandleFunc("/404", http.NotFoundHandler().ServeHTTP)
	limiter := rate.NewLimiter(100, 30)
	http.HandleFunc("/limited", handlerLimit(limiter, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%v<em>Hello, world</em>\n", htmlHead)))
	}))
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
