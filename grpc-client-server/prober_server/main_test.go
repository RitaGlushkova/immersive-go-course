package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"

	pb "github.com/RitaGlushkova/immersive-go-course/grpc-client-server/prober"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterProberServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type Test struct {
	req       *pb.ProbeRequest
	replyCode int64
}

func TestDoProbes(t *testing.T) {
	tests := map[string]Test{
		"success": {
			req: &pb.ProbeRequest{
				Endpoint:               "https://www.google.com/",
				NumberOfRequestsToMake: 3,
			},
			replyCode: 200,
		},
		"failed": {
			req: &pb.ProbeRequest{
				Endpoint:               "https://www.goog",
				NumberOfRequestsToMake: 1,
			},
			replyCode: 0,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			require.NoError(t, err)
			defer conn.Close()
			client := pb.NewProberClient(conn)
			resp, err := client.DoProbes(context.Background(), tt.req)
			require.NoError(t, err)
			for i := 0; i < int(tt.req.NumberOfRequestsToMake); i++ {
				fmt.Println(resp.Replies[i])
				require.Equal(t, tt.replyCode, resp.Replies[i].ReplyCode)
			}
		},
		)
	}
}

func XTestSomething(t *testing.T) {
	var largeNumber = 12345
	fmt.Println(largeNumber)

	var floatNumber = float64(largeNumber / 100)
	require.Equal(t, 123.45, floatNumber)

}

// test for a port to listen

// given the port is busy when we try to start it it will fail - return error

func TestSetPrometheusWhenPortBusy(t *testing.T) {
	// make port we are testing for busy
	lis, err := net.Listen("tcp", ":2112")
	require.NoError(t, err)
	defer lis.Close()
	// call our function, which we are testing
	err = setupPrometheus()
	// expect to return an error
	require.Error(t, err)
}

func TestSetPrometheusWhenPortAvail(t *testing.T) {

	// call our function, which we are testing
	err := setupPrometheus()
	// expect to return an error
	require.NoError(t, err)

	//make request to the port
	res, err := http.Get("http://localhost:2112/metrics")
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	// assert that some metrics are returner
}
