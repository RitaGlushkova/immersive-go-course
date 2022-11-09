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
	//setting up flags as Int helps to check if flag is a string that can not be converted to int.
	mcrouterPort := flag.Int("mcrouter", 11211, "Parses Mcrouter port")

	//probably could use flag.Var but I think this way covers a bit more checks.
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

	//repeating, might need to move to a struct
	mcSet := memcache.New("localhost:" + strconv.Itoa(*mcrouterPort))

	//harcoding value and key. Might not be best approach but it is a test?
	cacheErr := mcSet.Set(&memcache.Item{Key: "mykey", Value: []byte("4"), Expiration: 60})
	if cacheErr != nil {
		log.Fatal(cacheErr)
	}

	//collecting values and cache miss count in overview
	var overview Overview
	for _, port := range memcachePorts {
		mcGet := memcache.New("localhost:" + strconv.Itoa(port))
		val, err := mcGet.Get("mykey")

		// checking is err is not what we expected
		if err != nil && err.Error() != "memcache: cache miss" {
			log.Fatal(err)
		}

		// counting miss cache
		if err != nil && err.Error() == "memcache: cache miss" {
			overview.errCount++

			// collecting values from cache. I want to print it so see if it is what we want
		} else {
			overview.value = append(overview.value, string(val.Value))
		}
	}

	//Printing result and also a potential unexpected set up
	result := CheckCacheMode(overview, memcachePorts)
	log.Println(result)
}

func CheckCacheMode(overview Overview, nodes []int) string {
	if overview.errCount == (len(nodes)-1) && len(overview.value) == 1 {
		return fmt.Sprintf("caches are operating in sharded mode. Number of cache miss: %d, Values from cache: %v\n", overview.errCount, overview.value)
	} else if overview.errCount == 0 && len(overview.value) == len(nodes) {
		return fmt.Sprintf("caches are operating in replicated mode. Number of cache miss: %d, Values from cache: %v\n", overview.errCount, overview.value)
	} else {
		return fmt.Sprintf("caches are not operating correctly. Number of cache miss: %d. Values from cache: %v\n", overview.errCount, overview.value)
	}
}
