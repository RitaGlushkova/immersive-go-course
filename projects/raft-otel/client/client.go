// Package main implements a client for Prober service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	cmd "github.com/RitaGlushkova/raft-otel/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ports = flag.String("ports", "60051,60052,60053,60054,60055", "the addresses of servers to connect to")
	value = flag.Int64("value", 10, "dfines the size of a loop for generation values")
	key   = flag.String("key", "x", "number of requests to make")
)

func main() {
	flag.Parse()
	addresses := strings.Split(*ports, ",")
	for _, addr := range addresses {
		// Set up a connection to the server.
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := cmd.NewCommandClient(conn)

		// Contact the server and print out its response.
		resp, err := SendValue(c, &cmd.CommandRequest{Entry: &cmd.Entry{Key: *key, Value: *value}})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Response: %v", resp)
	}
}

func SendValue(c cmd.CommandClient, req *cmd.CommandRequest) (*cmd.CommandReply, error) {
	duration := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	resp, err := c.Store(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
