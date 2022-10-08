package static

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func Run() {

	path := flag.String("path", "assets", "files path")
	port := flag.String("port", "8082", "port listening")
	flag.Parse()
	log.Printf("path: %s\n", *path)
	log.Printf("port: %s\n", *port)
	fileServer := http.FileServer(http.Dir("./" + *path))
	http.Handle("/", fileServer)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), nil))
}

