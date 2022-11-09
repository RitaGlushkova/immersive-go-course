package main

import (
	//"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
)

type Overview struct {
	errCount int
	value    []string
}

func main() {
	mcrouterPort := flag.Int("mcrouter", 11211, "Parses Mcrouter port")
	var memcachePorts []int
	flag.Func("memcacheds", "Parses ports for memcache", func(flagValue string) error {
		stringsSlice := strings.Split(flagValue, ",")
		for _, str := range stringsSlice {
			port, err := strconv.Atoi(str)
			if err != nil {
				return fmt.Errorf("couldn't convert provided flag argument: %s into an integer. Error: %v", str, err)
			}
			memcachePorts = append(memcachePorts, port)
		}
		return nil
	})
	flag.Parse()
	mcSet := memcache.New("localhost:" + strconv.Itoa(*mcrouterPort))
	cacheErr := mcSet.Set(&memcache.Item{Key: "mykey", Value: []byte("4"), Expiration: 60})
	if cacheErr != nil {
		log.Fatal(cacheErr)
	}
	var overview Overview
	for _, port := range memcachePorts {
		mcGet := memcache.New("localhost:" + strconv.Itoa(port))
		val, err := mcGet.Get("mykey")
		if err != nil && err.Error() != "memcache: cache miss" {
			log.Fatal(err)
		}
		if err != nil && err.Error() == "memcache: cache miss" {
			overview.errCount++
		} else {
			overview.value = append(overview.value, string(val.Value))
		}
	}
	if overview.errCount == (len(memcachePorts)-1) && len(overview.value) == 1 {
		fmt.Printf("caches are operating in sharded mode. Number of cache miss: %d, Values from cache: %v\n", overview.errCount, overview.value)
	} else if overview.errCount == 0 && len(overview.value) == len(memcachePorts) {
		fmt.Printf("caches are operating in replicated mode. Number of cache miss: %d, Values from cache: %v\n", overview.errCount, overview.value)
	} else {
		fmt.Printf("caches are not operating correctly. Number of cache miss: %d. Values from cache: %v\n", overview.errCount, overview.value)
	}
}
