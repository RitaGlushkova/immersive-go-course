// Package main implements a client for Prober service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	cmd "github.com/RitaGlushkova/raft-otel/client_rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ports = flag.String("ports", "50051,50052,50053,50054,50055", "the addresses of servers to connect to")
	value = flag.Int64("value", 10, "dfines the size of a loop for generation values")
	key   = flag.String("key", "x", "number of requests to make")
)

func main() {
	flag.Parse()
	addresses := strings.Split(*ports, ",")
	fmt.Println(addresses)
	for i := 0; i < len(addresses); i++ {
		fmt.Printf("localhost:%v", addresses[i])
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", addresses[i]), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := cmd.NewCommandClient(conn)
		fmt.Println("CONNECTED to ", addresses[i])
		// Contact the server and print out its response.
		resp, err := SendValue(c, &cmd.Request{Entry: &cmd.Entry{Key: *key, Value: *value}})
		if err != nil {
			fmt.Println("ERROR SENDING VALUE")
			log.Fatal(err)
		}
		if resp.IsLeader {
			log.Printf("Response from server: %v", resp)
			break
		}
	}
}

func SendValue(c cmd.CommandClient, req *cmd.Request) (*cmd.Reply, error) {
	duration := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	resp, err := c.Store(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
