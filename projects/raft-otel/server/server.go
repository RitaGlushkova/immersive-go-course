package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	cmd "github.com/RitaGlushkova/raft-otel/command"
	"google.golang.org/grpc"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	leaderId = flag.Int64("leaderId", 1, "leader id")
)

type server struct {
	cmd.UnimplementedCommandServer
}
type Entry struct {
	Key   string
	Value int64
}

func (s *server) Store(ctx context.Context, in *cmd.CommandRequest) (*cmd.CommandReply, error) {
	clientAssumedByClient := in.GetLeaderId()
	if clientAssumedByClient != *leaderId {
		//return reply to client saying that you are not a leader
	}

	//if you are a leader, then do the job

	var entries = make([]*cmd.AcceptedEntry, 0)

	for _, entry := range in.GetEntries() {
		entries = append(entries, &cmd.AcceptedEntry{Key: entry.Key, Value: entry.Value})
	}
	return &cmd.CommandReply{AcceptedEntry: entries}, nil
}

func setupPrometheus(port int) (int, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, err
	}
	go http.Serve(lis, nil)
	return lis.Addr().(*net.TCPAddr).Port, nil
}

func main() {
	_, err := setupPrometheus(2112)
	if err != nil {
		log.Fatal("Failed to listen on port :2112", err)
	}
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	cmd.RegisterCommandServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
