package static

import (
	"flag"
	"fmt"
	"net/http"
)

func Run() {

	var path string
	var port string

	flag.StringVar(&path, "path", "", "files path")
	flag.StringVar(&port, "port", "", "port")
	flag.Parse()
	muxer := http.NewServeMux()
	if path != "" && port != "" {
		fileServer := http.FileServer(http.Dir(path))
		muxer.Handle("/", fileServer)
		http.ListenAndServe(fmt.Sprintf(":%s", port), muxer)
	}

}
