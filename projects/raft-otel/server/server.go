package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	cmd "github.com/RitaGlushkova/raft-otel/client_rpc"
	rt "github.com/RitaGlushkova/raft-otel/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	portC     = flag.Int64("portC", 50051, "The server port for client")
	peerPorts = flag.String("pp", "60052,60053,60054,60055", "peers ports")
	portR     = flag.Int64("portR", 60051, "server port for raft")
	//temporary flags
	leaderId = flag.Int64("leaderId", 60051, "leader id")
)

type Entry struct {
	Key   string
	Value int64
	Term  int64
}

type RaftServer struct {
	cmd.UnimplementedCommandServer
	rt.UnimplementedRaftServer

	leaderId int64
	//PersistentState
	currentTerm int64
	//votedFor    int64
	log []Entry
	// VolatileState
	commitIndex int64
	lastApplied int64

	// VolatileStateLeader
	//nextIndex  []int64
	//matchIndex []int64
}

func (s *RaftServer) Store(ctx context.Context, in *cmd.Request) (*cmd.Reply, error) {
	if *portR != s.leaderId {
		fmt.Printf("I am not a leader, i think the leaders id is %v", s.leaderId)
		return &cmd.Reply{NotLeaderMessage: "I am not a leader", IsLeader: false, SuggestedLeaderID: s.leaderId}, nil
	}
	clientEntry := in.GetEntry()
	fmt.Println("ENTRY FROM THE CLIENT", clientEntry)
	//Append to its log first
	s.log = append(s.log, Entry{Key: clientEntry.Key, Value: clientEntry.Value, Term: s.currentTerm})
	//Update lastApplied
	s.lastApplied = int64(len(s.log) - 1)

	//Act as a client to other servers
	addresses := strings.Split(*peerPorts, ",")
	successCounter := 0
	respChan := make(chan *rt.ResultAppend, len(addresses))
	//respLogs := make(chan string, len(addresses))
	for _, addr := range addresses {
		// Set up a connection to the servers in module.
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("couldn't connect to a peer with address %v, %v", addr, err)
		}
		defer conn.Close()
		raft := rt.NewRaftClient(conn)

		go func() {
			for {
				resp, err := TalkToPeers(raft, &rt.RequestAppend{Entry: &rt.Entry{Key: clientEntry.Key, Value: clientEntry.Value, Term: s.currentTerm}, LeaderId: *leaderId, PrevLogIndex: s.lastApplied, PrevLogTerm: s.log[s.lastApplied].Term, LeaderCommit: s.commitIndex, Term: s.currentTerm})
				if err != nil {
					fmt.Println(err)
				}
				if resp.Succeeds {
					respChan <- resp
					break
				}
			}
		}()
	}
	for i := 0; i < len(addresses); i++ {
		resp := <-respChan
		if resp.Succeeds {
			successCounter++
		}
		if successCounter == len(addresses)/2 {
			//+1 is yourself
			s.commitIndex = s.lastApplied
			//send heartbeat in go routine
			fmt.Println("COMMIT INDEX", s.commitIndex)
			fmt.Println("Last commited value", s.log[s.commitIndex].Value, "with key", s.log[s.commitIndex].Key)
			return &cmd.Reply{SuccessMessage: "Value stored", IsLeader: true}, nil
		}
	}

	return nil, fmt.Errorf("could not store value, try again")
}

// AppendValue requests to appends a value to the raft log.
func TalkToPeers(raft rt.RaftClient, req *rt.RequestAppend) (*rt.ResultAppend, error) {
	duration := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	resp, err := raft.AppendEntry(ctx, req)
	if err != nil {
		return resp, err
	}
	fmt.Println("I think I am a leader", req.LeaderId, "my current term: ", req.Term)
	return resp, nil
}

//func SendAHeartbeat

func (s *RaftServer) AppendEntry(ctx context.Context, in *rt.RequestAppend) (*rt.ResultAppend, error) {
	//append to log
	if in.Entry != nil {
		//check leader's term
		if in.Entry.Term < s.currentTerm {
			return &rt.ResultAppend{Succeeds: false, Term: s.currentTerm}, fmt.Errorf("i don't think you are a leader")
		}
		//check if log is consistent
		// if s.lastApplied != in.PrevLogIndex || s.log[s.lastApplied].Term != in.PrevLogTerm {
		// 	return &rt.ResultAppend{Succeeds: false, Term: s.currentTerm}, fmt.Errorf("log is not consistent")
		// }
		s.log = append(s.log, Entry{Key: in.Entry.Key, Value: in.Entry.Value, Term: s.currentTerm})
		fmt.Println("I am a follower, I have appended to my log", in.Entry.Key, in.Entry.Value, "with term", s.currentTerm, "my logs are", s.log)
		return &rt.ResultAppend{Succeeds: true, Term: s.currentTerm}, nil
	}
	// heartbeat
	if in.Term < s.currentTerm {
		return &rt.ResultAppend{Succeeds: false, Term: s.currentTerm}, fmt.Errorf("i don't think you are a leader")
	}
	//is current term is less or equal to leader's term
	s.currentTerm = in.Term
	if s.lastApplied != in.PrevLogIndex || s.log[s.lastApplied].Term != in.PrevLogTerm {
		return &rt.ResultAppend{Succeeds: false, Term: s.currentTerm}, fmt.Errorf("log is not consistent")
	}
	//Update lastApplied
	if s.commitIndex < in.LeaderCommit {
		s.commitIndex = in.LeaderCommit
		s.lastApplied = int64(len(s.log) - 1)
	}
	return &rt.ResultAppend{Succeeds: true, Term: s.currentTerm}, nil
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

	server := &RaftServer{leaderId: *leaderId, currentTerm: 0, commitIndex: 0, lastApplied: 0}
	server.leaderId = *leaderId
	fmt.Println("I just started. I think leader is ", server.leaderId, "my current logs: ", server.log, "my current term: ", server.currentTerm)
	//server for raft
	serverForRaft := grpc.NewServer()
	rt.RegisterRaftServer(serverForRaft, server)
	go func() {
		log.Printf("server listening at %v", lisRaft.Addr())
		if err := serverForRaft.Serve(lisRaft); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	//server for client
	serverForClient := grpc.NewServer()
	cmd.RegisterCommandServer(serverForClient, server)
	log.Printf("server listening at %v", lisClient.Addr())
	if err := serverForClient.Serve(lisClient); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
