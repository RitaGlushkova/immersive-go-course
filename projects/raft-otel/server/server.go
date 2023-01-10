package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/RitaGlushkova/raft-otel/prober"
	"google.golang.org/grpc"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	leaderId = flag.Int64("leaderId", 1, "leader id")
)

type server struct {
	pb.UnimplementedProberServer
}

type Entry struct {
	Key   string
	Value int64
}

func (s *server) DoProbes(ctx context.Context, in *pb.ProbeRequest) (*pb.ProbeReply, error) {
	clientAssumedByClient := in.GetLeaderId()
	if clientAssumedByClient != *leaderId {
		//return reply to client saying that you are not a leader
	}

	//if you are a leader, then do the job

	var entries = make([]*pb.AcceptedEntry, 0)

	for _, entry := range in.GetEntries() {
		entries = append(entries, &pb.AcceptedEntry{Key: entry.Key, Value: entry.Value})
	}
	return &pb.ProbeReply{AcceptedEntry: entries}, nil
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
	pb.RegisterProberServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
