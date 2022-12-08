package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			for {
				fmt.Println(i)
				time.Sleep(1 * time.Second)
			}
		}()
		time.Sleep(1 * time.Second)
	}
	time.Sleep(60 * time.Second)
}
