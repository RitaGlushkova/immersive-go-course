package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	cmd "github.com/RitaGlushkova/raft-otel/command"
	rt "github.com/RitaGlushkova/raft-otel/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	portC     = flag.Int64("portC", 50051, "The server port for client")
	peerPorts = flag.String("pp", "60052,60053,60054,60055", "peers ports")
	portR     = flag.Int64("portR", 60051, "server port for raft")
	leader    = flag.Bool("leader", false, "Am I a leader?")
	leaderId  = flag.Int64("leaderId", 0, "leader id")
)

type Entry struct {
	Key   string
	Value int64
	//term  int64
}

type RaftServer struct {
	cmd.UnimplementedCommandServer
	rt.UnimplementedRaftServer

	//leaderId int64
	//PersistentState
	//currentTerm int64
	//votedFor    int64
	log []Entry
	// VolatileState
	commitIndex int64
	lastApplied int64

	// VolatileStateLeader
	nextIndex  []int64
	matchIndex []int64
}

func (s *RaftServer) Store(ctx context.Context, in *cmd.CommandRequest) (*cmd.CommandReply, error) {
	if !*leader {
		return &cmd.CommandReply{NotLeaderMessage: "I am not a leader", IsLeader: false}, nil
	}
	clientEntry := in.GetEntry()
	//Act as a client to other servers
	addresses := strings.Split(*peerPorts, ",")
	successCounter := 0
	respChan := make(chan *rt.ResultAppend, 1)
	for _, addr := range addresses {
		// Set up a connection to the servers in module.
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		raft := rt.NewRaftClient(conn)

		go func() {
			resp, err := TellPeerToAppend(raft, &rt.RequestAppend{Entry: &rt.Entry{Key: clientEntry.Key, Value: clientEntry.Value}, LeaderId: *leaderId, PrevLogIndex: int64(len(s.log) - 1)})
			if err != nil {
				fmt.Println(err)
			}
			respChan <- resp
		}()
	}
	for i := 0; i < len(addresses); i++ {
		resp := <-respChan
		if resp.Succeeds {
			successCounter++
		}
	}
	if successCounter > len(addresses)/2 {
		//run command to append it to the leader
		return &cmd.CommandReply{SuccessMessage: "Value stored"}, nil
	}
	return nil, fmt.Errorf("could not store value, try again")
}

// AppendValue requests to appends a value to the raft log.
func TellPeerToAppend(raft rt.RaftClient, req *rt.RequestAppend) (*rt.ResultAppend, error) {
	duration := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	resp, err := raft.AppendEntry(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *RaftServer) AppendEntry(ctx context.Context, in *rt.RequestAppend) (*rt.ResultAppend, error) {
	if in.LeaderId != *leaderId {
		return &rt.ResultAppend{Succeeds: false}, fmt.Errorf("i don't think you are a leader")
	}
	if in.PrevLogIndex != s.lastApplied {
		return &rt.ResultAppend{Succeeds: false}, fmt.Errorf("my previous log index is not the same as yours")
	}
	if s.log[len(s.log)-1].Key != in.Entry.Key || s.log[len(s.log)-1].Value != in.Entry.Value {
		return &rt.ResultAppend{Succeeds: false}, fmt.Errorf("my last log entry is not the same as yours")
	}
	s.log = append(s.log, Entry{Key: in.Entry.Key, Value: in.Entry.Value})
	return &rt.ResultAppend{Succeeds: true}, nil
}

func main() {
	flag.Parse()

	//Act as a server
	lisClient, err := net.Listen("tcp", fmt.Sprintf(":%d", *portC))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	lisRaft, err := net.Listen("tcp", fmt.Sprintf(":%d", *portR))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//server for raft
	serverForRaft := grpc.NewServer()
	rt.RegisterRaftServer(serverForRaft, &RaftServer{})
	log.Printf("server listening at %v", lisRaft.Addr())
	if err := serverForRaft.Serve(lisRaft); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//server for client
	serverForClient := grpc.NewServer()
	cmd.RegisterCommandServer(serverForClient, &RaftServer{})
	log.Printf("server listening at %v", lisClient.Addr())
	if err := serverForClient.Serve(lisClient); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
