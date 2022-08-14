package static

import (
	"log"
	"net/http"
	"path/filepath"
)

type Config struct {
	Dir string
}

func Run(config Config) {
	// The "/" path handles everything, so we need to inspect the path (`req.URL.Path`) to be able to
	// identify which file to serve.
	// https://pkg.go.dev/net/http#ServeMux.Handle
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Build a full absolute path to the file, relative to the config.Dir
		path := filepath.Join(config.Dir, r.URL.EscapedPath())
		log.Println(r.Method, r.URL.EscapedPath(), path)
		http.ServeFile(w, r, path)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}