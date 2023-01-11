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
	leaderId  = flag.Int64("leaderId", 1, "leader id")
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
	//nextIndex  []int64
	//matchIndex []int64
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
			fmt.Printf("couldn't connect to a peer with address %v, %v", addr, err)
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
		s.log = append(s.log, Entry{Key: clientEntry.Key, Value: clientEntry.Value})
		s.commitIndex = int64(len(s.log) - 1)
		s.lastApplied = int64(len(s.log) - 1)
		respChan := make(chan *rt.ResponseUpdate, 1)

		//tell other servers to update their commit index
		for _, addr := range addresses {
			// Set up a connection to the servers in module.
			conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				fmt.Printf("couldn't connect to a peer with address %v, %v", addr, err)
			}
			defer conn.Close()
			raft := rt.NewRaftClient(conn)
			go func() {
				resp, err := TellPeerToUpdateCommit(raft, &rt.RequestUpdate{LeaderCommit: s.commitIndex, LeaderId: *leaderId})
				if err != nil {
					fmt.Println(err)
				}
				respChan <- resp
			}()
			//repeat until get success HOW????
		}

		//run command to append it to the leader but run update in the background
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

func TellPeerToUpdateCommit(raft rt.RaftClient, req *rt.RequestUpdate) (*rt.ResponseUpdate, error) {
	duration := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	resp, err := raft.UpdateCommit(ctx, req)
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
	s.log = append(s.log, Entry{Key: in.Entry.Key, Value: in.Entry.Value})
	s.lastApplied = int64(len(s.log) - 1)
	return &rt.ResultAppend{Succeeds: true}, nil
}

func UpdateCommit(s *RaftServer, in *rt.RequestUpdate) (*rt.RequestUpdate, error) {
	if in.LeaderId != *leaderId {
		return &rt.RequestUpdate{LeaderCommit: s.commitIndex}, fmt.Errorf("i don't think you are a leader")
	}
	s.commitIndex = in.LeaderCommit
	return &rt.RequestUpdate{LeaderCommit: in.LeaderCommit}, nil
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
