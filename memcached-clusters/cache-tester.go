package main

import (
	"fmt"
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

type Overview struct {
	err string
	val string
}

func main() {
	mcSet := memcache.New("localhost:11211")
	cacheErr := mcSet.Set(&memcache.Item{Key: "mykey", Value: []byte("4"), Expiration: 60})
	if cacheErr != nil {
		log.Fatal(cacheErr)
	}

	var memcacheds = []string{"11212", "11213", "11214"}
	overviewSlice := make([]Overview, 0)
	for _, memcached := range memcacheds {
		mcGet := memcache.New(fmt.Sprintf("localhost:%v", memcached))
		val, err := mcGet.Get("mykey")
		if err != nil {
			overviewSlice = append(overviewSlice, Overview{err.Error(), ""})
		} else {
			overviewSlice = append(overviewSlice, Overview{"", string(val.Value)})
		}

	}
	fmt.Println(overviewSlice)
}
