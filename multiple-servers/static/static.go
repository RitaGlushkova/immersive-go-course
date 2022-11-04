package static

import (
	"fmt"
	"net/http"
)

type Config struct{
	Port int
	Path string
}

func Run(c Config) error {
	fileServer := http.FileServer(http.Dir(c.Path))
	http.Handle("/", fileServer)
	return http.ListenAndServe(fmt.Sprintf(":%v", c.Port), nil)
}

