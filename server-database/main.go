package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func encodedMarshalJSON(data interface{}, queryVal string) ([]byte, error) {
	indent, errIndent := strconv.Atoi(queryVal)
	var marshalData []byte
	var marshalErr error
	if errIndent != nil {
		//DO I WANT TO INFORM ABOUT IT ????
	}
	if indent > 0 && errIndent == nil {
		marshalData, marshalErr = json.MarshalIndent(data, "", strings.Repeat(" ", indent))
	} else {
		marshalData, marshalErr = json.Marshal(data)
	}
	if marshalErr != nil {
		fmt.Fprintf(os.Stderr, "Couldn't encode inserted values: %v\n", marshalErr)
	}
	return marshalData, marshalErr
}

func fetchImages(conn *pgx.Conn) ([]Image, error) {
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

func postImage(conn *pgx.Conn, r *http.Request) (*Image, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read body: %v\n", err)
		return nil, err
	}
	var img Image
	err = json.Unmarshal(b, &img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't decode input: %v\n", err)
		return nil, err
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO public.images (title, url, alt_text) VALUES ($1, $2, $3)", img.Title, img.URL, img.AltText)
	if err != nil {
		fmt.Fprintf(os.Stderr, "POST query failed: %v\n", err)
		return nil, err
	}
	return &img, nil
}

func handlerImages(w http.ResponseWriter, r *http.Request) {
	//not sure it is a good solution
	_, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		fmt.Fprintf(os.Stderr, "Environment variable is not set")
		os.Exit(1)
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	queryVal := r.URL.Query().Get("indent")

	switch r.Method {
	case "GET":
		images, err := fetchImages(conn)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print(err)
		}
		encoded, err := encodedMarshalJSON(images, queryVal)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(encoded))
	case "POST":
		img, err := postImage(conn, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print(err)
		}
		encoded, err := encodedMarshalJSON(img, queryVal)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(encoded))
	default:
		w.WriteHeader(405)
		w.Write([]byte("Only GET and POST methods are available"))

	}
}

func main() {

	http.HandleFunc("/images.json", handlerImages)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
