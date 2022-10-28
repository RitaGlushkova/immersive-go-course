package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"

	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

type Image struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	AltText string `json:"alt_text"`
	Pixels  int    `json:"pixels"`
}

type Server struct {
	conn *pgx.Conn
}

func main() {
	godotenv.Load()
	env := os.Getenv("DATABASE_URL")
	if env == "" {
		fmt.Fprintf(os.Stderr, "Environment variable is not set\n")
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
		imagesJsonInBytes, err := MarshalJSON(images, queryVal, os.Stderr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(imagesJsonInBytes)
	case "POST":
		img, err := saveImage(s.conn, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		imageJsonInBytes, err := MarshalJSON(img, queryVal, os.Stderr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(imageJsonInBytes)
	default:
		w.WriteHeader(405)
		w.Write([]byte("Only GET and POST methods are available"))

	}
}

// turns our struct into JSON
func MarshalJSON(data interface{}, queryVal string, diagnostics io.Writer) ([]byte, error) {
	indent, errIndent := strconv.Atoi(queryVal)
	var marshalData []byte
	var marshalErr error
	if errIndent != nil {
		fmt.Println(errIndent, "default indent will be set to 0")
	}
	if indent > 0 && indent < 15 {
		marshalData, marshalErr = json.MarshalIndent(data, "", strings.Repeat(" ", indent)) //returns Json encoded value in []byte with indentation
	} else {
		marshalData, marshalErr = json.Marshal(data) //returns Json encoded value in []byte
	}
	if marshalErr != nil {
		fmt.Fprintf(diagnostics, "couldn't proceed with Marshal: %v\n", marshalErr)
		return nil, marshalErr
	}
	return marshalData, nil
}

func FetchImages(conn *pgx.Conn) ([]Image, error) {
	rows, err := conn.Query(context.Background(), "SELECT title, url, alt_text, pixels FROM public.images")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetching query failed: %v\n", err)
		return nil, err
	}
	var images []Image
	for rows.Next() {
		var title, url, altText string
		var pixels int
		err = rows.Scan(&title, &url, &altText, &pixels)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scan for rows failed: %v\n", err)
			return nil, err
		}
		images = append(images, Image{Title: title, URL: url, AltText: altText, Pixels: pixels})
	}
	return images, nil
}

func saveImage(conn *pgx.Conn, body io.Reader) (*Image, error) {
	//NewDecoder returns a new decoder that reads from r.
	//The decoder introduces its own buffering and may read data from r beyond the JSON values requested.
	decoder := json.NewDecoder(body)
	var img Image
	var exists bool
	err := decoder.Decode(&img) //decodes and stores json.encoded value into var img (Image struct)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't decode input: %v\n", err)
		return nil, err
	}

	//check for valid URL
	resp, err := http.Get(img.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't process with Get request %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	//getting image resolution
	config, err := jpeg.DecodeConfig(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't get image resolution, image won't be saved: %v", err)
	}
	// The resolution is height x width
	img.Pixels = config.Height * config.Width
	//check if URL exists already
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid URL link image won't be saved: %v", err)
	}
	err = conn.QueryRow(context.Background(), "SELECT exists (SELECT url FROM public.images WHERE url = $1)", img.URL).Scan(&exists)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error querying: %v\n", err)
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("image with URL: %v already exists", img.URL)
	}

	//Check for Alt text
	if img.AltText == "" || compareTitleAndAltText(img.AltText, img.Title) < 1 {
		return nil, fmt.Errorf("AltText is not descriptive enough: %s. Image won't be saved", img.AltText)
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO public.images (title, url, alt_text, pixels) VALUES ($1, $2, $3, $4)", img.Title, img.URL, img.AltText, img.Pixels)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't insert row %v\n", err)
		return nil, err
	}
	return &img, nil
}

func compareTitleAndAltText(altText, title string) int {
	a := strings.Split(altText, " ")
	count := 0
	for _, v := range a {
		if strings.Contains(strings.ToLower(title), strings.ToLower(v)) {
			count++
		}
	}
	return count
}
