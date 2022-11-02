package main

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

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
	pb.RegisterGreeterServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSayHello(t *testing.T) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewProberClient(conn)
	resp, err := client.DoProbes(context.Background(), &pb.ProbeRequest{Endpoint: "https://www.google.com", NumberOfRequestsToMake: 6})
	require.NoError(t, err)
	require.Equal(t, 3*time.Second, resp.AverageLatencyMsecs.AsDuration())
}
