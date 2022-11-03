// Package main implements a client for Prober service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/CodeYourFuture/immersive-go-course/grpc-client-server/prober"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr        = flag.String("addr", "localhost:50051", "the address to connect to")
	endpoint    = flag.String("endpoint", "http://www.google.com", "defines endpoint")
	numberOfReq = flag.Int64("n", 1, "number of requests to make")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProberClient(conn)

	// Contact the server and print out its response.
	ProbeLog(c, &pb.ProbeRequest{Endpoint: *endpoint, NumberOfRequestsToMake: *numberOfReq})
}

func ProbeLog(c pb.ProberClient, req *pb.ProbeRequest) *pb.ProbeReply {
	duration := 1 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	resp, err := c.DoProbes(ctx, req)
	if err != nil {
		log.Fatalf("could not probe: %v", err)
		return nil
	}
	log.Printf("Average Latency for %d request(s) is %v milliseconds. %v", *numberOfReq, resp.GetAverageLatencyMsecs(), resp.Replies)
	return resp
}
