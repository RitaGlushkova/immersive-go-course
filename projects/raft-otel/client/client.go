// Package main implements a client for Prober service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	cmd "github.com/RitaGlushkova/raft-otel/command"
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
	c := cmd.NewCommandClient(conn)

	// Contact the server and print out its response.
	resp, err := SendValue(c, &cmd.CommandRequest{Entries: []*cmd.Entry{{Key: *key, Value: *value}}, LeaderId: *leaderId})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %v", resp)
}

func SendValue(c cmd.CommandClient, req *cmd.CommandRequest) (*cmd.CommandReply, error) {
	duration := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	resp, err := c.Store(ctx, req)
	if err != nil {
		return nil, status.Errorf(13, "could not store value: %v", err)
	}
	return resp, nil
}
