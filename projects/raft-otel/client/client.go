// Package main implements a client for Prober service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/RitaGlushkova/raft-otel/prober"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	addr     = flag.String("addr", "localhost:50051", "the address to connect to")
	value    = flag.Int64("value", 10, "dfines the size of a loop for generation values")
	key      = flag.String("key", "x", "number of requests to make")
	leaderId = flag.Int64("leaderId", 1, "leader id")
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
	_, err = ProbeLog(c, &pb.ProbeRequest{Entries: []*pb.Entry{{Key: *key, Value: *value}}, LeaderId: *leaderId})
	if err != nil {
		log.Fatal(err)
	}
}

func ProbeLog(c pb.ProberClient, req *pb.ProbeRequest) (*pb.ProbeReply, error) {
	duration := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	resp, err := c.DoProbes(ctx, req)
	if err != nil {
		return nil, status.Errorf(13, "could not probe: %v", err)
	}
	//log.Printf("Average Latency for %d request(s) is %v milliseconds. %v", *numberOfReq, resp.GetAverageLatencyMsecs(), resp.Replies)
	return resp, nil
}
