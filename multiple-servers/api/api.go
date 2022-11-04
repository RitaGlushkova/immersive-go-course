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
	
	return http.ListenAndServe(fmt.Sprintf(":%d",c.Port), nil)
}

func (s *Server) handlerImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	indent := r.URL.Query().Get("indent")
	switch r.Method {
	case "GET":
		images, err := FetchImages(s.conn)
		if err != nil {
			errMessage := "Something went wrong"
			http.Error(w, errMessage, http.StatusInternalServerError)
			return
		}
		encodeAndResponseJSON(&w, images, indent)

	case "POST":
		img, err := saveImage(s.conn, r.Body)
		if err != nil {
			errMessage := "Couldn't post an image"
			http.Error(w, errMessage, http.StatusBadRequest)
			return
		}
		encodeAndResponseJSON(&w, img, indent)

	default:
		w.WriteHeader(405)
		w.Write([]byte("Only GET and POST methods are available"))
	}
	log.Println(r.Method, r.URL.EscapedPath())
}
