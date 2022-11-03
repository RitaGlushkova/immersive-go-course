package main

import (
	"context"
	"log"
	"net"
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
				require.Equal(t, tt.replyCode, resp.Replies[i].ReplyCode)
			}
		},
		)
	}
}
