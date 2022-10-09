package api

import (
	"net/http"
	"os"
)

func encodeAndResponseJSON (w *http.ResponseWriter, data interface {}, query string) {
	encoded, err := EncodedMarshalJSON(data, query, os.Stderr)
		if err != nil {
			http.Error((*w), err.Error(), http.StatusInternalServerError)
			return

		}
		(*w).Header().Set("Content-Type", "application/json")
		(*w).Write([]byte(encoded))
}