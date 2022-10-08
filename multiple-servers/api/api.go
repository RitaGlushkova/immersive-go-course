package api

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
)

type Image struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	AltText string `json:"alt_text"`
}

type Server struct {
	conn *pgx.Conn
}

func Run() {
	var port int

	flag.IntVar(&port, "port", 8081, "Port is listening")
	flag.Parse()
	log.Printf("port: %d\n", port)
	env := os.Getenv("DATABASE_URL")
	if env == "" {
		fmt.Fprintf(os.Stderr, "Environment variable is not set")
		os.Exit(1)
	}
	conn, err := pgx.Connect(context.Background(), env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	s := &Server{conn: conn}
	http.HandleFunc("/images.json", s.handlerImages)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func (s *Server) handlerImages(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	queryVal := r.URL.Query().Get("indent")
	switch r.Method {
	case "GET":
		images, err := FetchImages(s.conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		encoded, err := EncodedMarshalJSON(images, queryVal, os.Stderr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(encoded))
	case "POST":
		img, err := saveImage(s.conn, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		encoded, err := EncodedMarshalJSON(img, queryVal, os.Stderr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(encoded))
	default:
		w.WriteHeader(405)
		w.Write([]byte("Only GET and POST methods are available"))
	}
	log.Println(r.Method, r.URL.EscapedPath())
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8082")
}
