package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
)

type Image struct {
	Title   string 
	URL     string 
	AltText string 
}

type Server struct {
	conn *pgx.Conn
}

type Config struct {
	Port         int
	DATABASE_URL string
}

func Run(c Config) error {

	conn, err := pgx.Connect(context.Background(), c.DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	s := &Server{conn: conn}
	http.HandleFunc("/images.json", s.handlerImages)
	return http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil)
}

func (s *Server) handlerImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
