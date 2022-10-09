package main

import (
	//"net/http"
	"flag"
	"fmt"
	"log"
	"os"
	"servers/api"
)

func main() {
	port := flag.Int("port", 8081, "Port is listening")
	flag.Parse()
	log.Printf("port: %d\n", *port)
	env := os.Getenv("DATABASE_URL")
	if env == "" {
		fmt.Fprintf(os.Stderr, "Environment variable is not set")
		os.Exit(1)
	}
	log.Fatal(api.Run(api.Config{
		Port:         *port,
		DATABASE_URL: os.Getenv("DATABASE_URL"),
	}))
}
