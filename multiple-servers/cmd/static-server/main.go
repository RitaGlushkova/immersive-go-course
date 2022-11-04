package main

import (
	"flag"
	"log"
	"servers/static"
)

func main() {
	path := flag.String("path", "assets", "files path")
	port := flag.Int("port", 8082, "port listening")
	flag.Parse()
	log.Printf("path: %s\n", *path)
	log.Printf("port: %d\n", *port)
	log.Fatal(static.Run(static.Config{
		Port: *port,
		Path: *path,
	}))
}
