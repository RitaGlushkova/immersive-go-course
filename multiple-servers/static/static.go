package static

import (
	"flag"
	"fmt"
	"net/http"
)

func Run() {

	var path string
	var port int

	flag.StringVar(&path, "path", "assets", "files path")
	flag.IntVar(&port, "port", 8082, "port listening")
	flag.Parse()
	muxer := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(path))
	muxer.Handle("/", fileServer)
	http.ListenAndServe(fmt.Sprintf(":%v", port), muxer)

}
