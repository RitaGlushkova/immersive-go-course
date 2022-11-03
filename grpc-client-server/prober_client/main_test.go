package main

import (
	"bytes"
	"context"
	"log"
	"net"
	"testing"

	pb "github.com/CodeYourFuture/immersive-go-course/grpc-client-server/prober"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type MockProbeServer struct {
	pb.UnimplementedProberServer
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterProberServer(s, &MockProbeServer{})
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
			var buf bytes.Buffer
			log.SetOutput(&buf)
			require.NoError(t, err)
			resp, err := ProbeLog(client, tt.req)
			if name == "success" {
				want := "Average Latency for 1 request(s) is 6 milliseconds. [latency_msecs:6  reply_code:200]"
				got := buf.String()
				require.Contains(t, got, want)
				require.NoError(t, err)
				require.Equal(t, float32(6), resp.AverageLatencyMsecs)
			}
			if name == "failed" {
				got := buf.String()
				want := `error_message:"Error: Get \"https://www.goog/\""`
				require.Contains(t, got, want)
			}
			require.Equal(t, tt.replyCode, resp.Replies[0].ReplyCode)
		},
		)
	}
}

func (*MockProbeServer) DoProbes(ctx context.Context, req *pb.ProbeRequest) (*pb.ProbeReply, error) {
	if req.Endpoint == "https://www.goog" {
		return &pb.ProbeReply{
			AverageLatencyMsecs: float32(3),
			Replies: []*pb.Reply{{
				LatencyMsecs: float32(3),
				ErrorMessage: "Error: Get \"https://www.goog/\"",
				ReplyCode:    0,
			}}}, nil
	} else {
		return &pb.ProbeReply{
			AverageLatencyMsecs: float32(6),
			Replies: []*pb.Reply{{
				LatencyMsecs: float32(6),
				ErrorMessage: "",
				ReplyCode:    200,
			}}}, nil
	}
}
