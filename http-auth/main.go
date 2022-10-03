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

type Server struct {
	username string
	password string
	limiter  *rate.Limiter
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		w.Header().Set("Content-Type", "text/html") //text/plain sends a string as response
		defaultResponse := htmlHead + "<em>Hello, world</em>"
		w.Write([]byte(fmt.Sprintf("%v\n", defaultResponse)))

		if len(r.URL.Query()) != 0 {
			w.Write([]byte("<p>Query parameters:<ul>"))

			for key, values := range r.URL.Query() {
				for _, value := range values {
					w.Write([]byte(fmt.Sprintf("<li>%v: [%v]</li>", key, html.EscapeString(value))))
				}

			}
			w.Write([]byte("</ul>"))
		}

	case "POST":
		w.Header().Set("Content-Type", "text/html")
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print("Error reading body of POST to index",err)
			w.WriteHeader(500)
			return
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

func (s *Server) handlerAuth(delegate func(username string) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != s.username || p != s.password {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			w.WriteHeader(401)
			return
		}
		w.Write([]byte(delegate(s.username)))
	}
}

func (s *Server) handlerLimit(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Allow() {
		w.WriteHeader(429)
	} else {
		w.Write([]byte(fmt.Sprintf("%v<em>Hello, world</em>\n", htmlHead)))
	}
}

func main() {
	s := Server{
		username: goDotEnvVariable("AUTH_USERNAME"),
		password: goDotEnvVariable("AUTH_PASSWORD"),
		limiter:  rate.NewLimiter(100, 30),
	}

	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/200", handler200)

	http.HandleFunc("/authenticated", s.handlerAuth(func(username string) string {
		return fmt.Sprintf("%v<em>Hello, %s!</em>\n", htmlHead, username)
	}))
	http.HandleFunc("/500", handler500)
	http.HandleFunc("/404", http.NotFoundHandler().ServeHTTP)
	http.HandleFunc("/limited", s.handlerLimit)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
