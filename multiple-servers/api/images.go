package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jackc/pgx/v4"
)

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

func saveImage(conn *pgx.Conn, body io.Reader) (*Image, error) {
	decoder := json.NewDecoder(body)
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