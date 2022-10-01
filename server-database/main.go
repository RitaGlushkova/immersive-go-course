package main

import (
	"context"
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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

func main() {
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
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/images.json", s.handlerImages)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func handlerIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func (s *Server) handlerImages(w http.ResponseWriter, r *http.Request) {
	queryVal := r.URL.Query().Get("indent")
	switch r.Method {
	case "GET":
		images, err := FetchImages(s.conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		encoded, err := EncodedMarshalJSON(images, queryVal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(encoded))
	case "POST":
		img, err := saveImage(s.conn, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		encoded, err := EncodedMarshalJSON(img, queryVal)
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
}

func EncodedMarshalJSON(data interface{}, queryVal string) ([]byte, error) {
	indent, errIndent := strconv.Atoi(queryVal)
	var marshalData []byte
	var marshalErr error
	if errIndent != nil {
		//DO I WANT TO INFORM ABOUT IT ????
	}
	if indent > 0 && indent < 15 && errIndent == nil {
		marshalData, marshalErr = json.MarshalIndent(data, "", strings.Repeat(" ", indent))
	} else {
		marshalData, marshalErr = json.Marshal(data)
	}
	if marshalErr != nil {
		fmt.Fprintf(os.Stderr, "Couldn't encode inserted values: %v\n", marshalData)
		return nil, marshalErr
	}
	return marshalData, nil
}

func FetchImages(conn *pgx.Conn) ([]Image, error) {
	rows, err := conn.Query(context.Background(), "SELECT title, url, alt_text FROM public.images")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		return nil, err
	}
	var images []Image
	for rows.Next() {
		var title, url, altText string
		err = rows.Scan(&title, &url, &altText)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Scan for rows failed: %v\n", err)
			return nil, err
		}
		images = append(images, Image{Title: title, URL: url, AltText: altText})
	}
	return images, nil
}

func saveImage(conn *pgx.Conn, r *http.Request) (*Image, error) {
	decoder := json.NewDecoder(r.Body)
	var img Image
	err := decoder.Decode(&img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't decode input: %v\n", err)
		return nil, err
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO public.images (title, url, alt_text) VALUES ($1, $2, $3)", img.Title, img.URL, img.AltText)
	if err != nil {
		return nil, fmt.Errorf("could not insert row: %w", err)
	}
	return &img, nil
}

