package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	//"net/http"

	cmd "github.com/RitaGlushkova/raft-otel/command"
	rt "github.com/RitaGlushkova/raft-otel/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	leaderId = flag.Int64("leaderId", 1, "leader id")
)

func (s *ServerClient) Store(ctx context.Context, in *cmd.CommandRequest) (*cmd.CommandReply, error) {
	clientAssumedByClient := in.GetLeaderId()
	if clientAssumedByClient != *leaderId {
		return &cmd.CommandReply{NotLeaderMessage: "I am not a leader", IsLeader: false}, nil
	}

	//if you are a leader, then do the job

	var entries = make([]*cmd.AcceptedEntry, 0)

	for _, entry := range in.GetEntries() {
		entries = append(entries, &cmd.AcceptedEntry{Key: entry.Key, Value: entry.Value})
	}
	return &cmd.CommandReply{AcceptedEntry: entries, IsLeader: true}, nil
}

func (s *ServerRaft) AppendEntries(ctx context.Context, in *rt.RequestAppend) (*rt.ResultAppend, error) {
	return &rt.ResponseAppend{Term: 1, Success: true}, nil
}

// func setupPrometheus(port int) (int, error) {
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
// 	if err != nil {
// 		return 0, err
// 	}
// 	go http.Serve(lis, nil)
// 	return lis.Addr().(*net.TCPAddr).Port, nil
// }

func main() {
	// _, err := setupPrometheus(2112)
	// if err != nil {
	// 	log.Fatal("Failed to listen on port :2112", err)
	// }
	flag.Parse()

	//Act as a client to other servers
	conn, err := grpc.Dial(fmt.Sprintf(":%d", *port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	raft := rt.NewRaftClient(conn)

	resp, err := AppendValue(raft, &rt.RequestAppend{Term: 1, LeaderId: 1, Entries: []*rt.Entry{{Key: "key", Value: 10}}})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %v", resp)

	//Act as a server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//server for raft
	serverForRaft := grpc.NewServer()
	rt.RegisterRaftServer(serverForRaft, &serverRaft{})
	log.Printf("server listening at %v", lis.Addr())
	if err := serverForRaft.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//server for client
	serverForClient := grpc.NewServer()
	cmd.RegisterCommandServer(serverForClient, &serverClient{})
	log.Printf("server listening at %v", lis.Addr())
	if err := serverForClient.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// AppendValue requests to appends a value to the raft log.
func AppendValue(raft rt.RaftClient, req *rt.RequestAppend) (*rt.ResultAppend, error) {
	duration := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	resp, err := raft.AppendEntries(ctx, req)
	if err != nil {
		return nil, status.Errorf(13, "could not store value: %v", err)
	}
	return resp, nil
}
