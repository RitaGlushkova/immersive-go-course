package main

import (
	"context"
	//"fmt"
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"

	pb "github.com/RitaGlushkova/immersive-go-course/grpc-client-server/prober"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func setNewServer(ctx context.Context) (pb.ProberClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pb.RegisterProberServer(baseServer, &server{})
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()
	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}
	client := pb.NewProberClient(conn)
	return client, closer
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
			ctx := context.Background()
			client, closer := setNewServer(ctx)
			defer closer()
			resp, err := client.DoProbes(context.Background(), tt.req)
			require.NoError(t, err)
			for i := 0; i < int(tt.req.NumberOfRequestsToMake); i++ {
				require.Equal(t, tt.replyCode, resp.Replies[i].ReplyCode)
			}
		},
		)
	}
}

func TestSomething(t *testing.T) {
	var largeNumber = float64(12345)
	fmt.Println(largeNumber)

	var floatNumber = largeNumber / 100
	fmt.Println(floatNumber)
	require.Equal(t, 123.45, floatNumber)

}

func TestSetPrometheusWhenPortBusy(t *testing.T) {
	// make port we are testing for busy
	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer lis.Close()
	// call our function, which we are testing
	_, err = setupPrometheus(lis.Addr().(*net.TCPAddr).Port)
	// expect to return an error
	require.Error(t, err)
}

func TestSetPrometheusWhenPortAvail(t *testing.T) {

	// call our function, which we are testing
	port, err := setupPrometheus(0)
	// expect to return an error
	require.NoError(t, err)

	//make request to the port
	res, err := http.Get(fmt.Sprintf("http://localhost:%v/metrics", port))
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	// assert that some metrics are returner
}
