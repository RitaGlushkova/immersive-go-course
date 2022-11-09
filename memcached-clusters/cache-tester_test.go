package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Tests struct {
	mockOverview Overview
	nodes        []int
	result       string
}

func TestCheckCacheMode(t *testing.T) {

	tests := map[string]Tests{
		"Should return replicated mode": {
			mockOverview: Overview{
				errCount: 0,
				value:    []string{"4", "4", "4"},
			},
			nodes:  []int{23, 24, 44},
			result: "caches are operating in replicated mode. Number of cache miss: 0, Values from cache: [4 4 4]\n",
		},
		"Should return sharded mode": {
			mockOverview: Overview{
				errCount: 2,
				value:    []string{"4"},
			},
			nodes:  []int{23, 24, 44},
			result: "caches are operating in sharded mode. Number of cache miss: 2, Values from cache: [4]\n",
		},
		"Should return incorrect set up mode": {
			mockOverview: Overview{
				errCount: 2,
				value:    []string{"4", "4"},
			},
			nodes:  []int{23, 24, 44},
			result: "caches are not operating correctly. Number of cache miss: 2. Values from cache: [4 4]\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := CheckCacheMode(tt.mockOverview, tt.nodes)
			require.Equal(t, tt.result, result)
		})
	}
}

//check for errors if flags are not parsed correctly???

//mock flags and test main????
