package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
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
	state    string
	leaderId int64
	//PersistentState
	currentTerm int64
	//votedFor    int64
	log []Entry
	// VolatileState
	commitIndex int64
	lastApplied int64
	peerPorts   []int64
	// VolatileStateLeader
	nextIndex []int64
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
	s.lastApplied = int64(len(s.log) - 1)
	fmt.Println("last applied on the leader", s.lastApplied)
	//Act as a client to other servers
	successCounter := 0
	respChan := make(chan *rt.ResultAppend, len(s.peerPorts))
	//respLogs := make(chan string, len(addresses))
	for _, port := range s.peerPorts {
		// Set up a connection to the servers in module.
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("couldn't connect to a peer with address %v, %v", port, err)
		}
		defer conn.Close()
		raft := rt.NewRaftClient(conn)
		go func() {
			for {
				resp, err := TalkToPeers(raft, &rt.RequestAppend{Entry: &rt.Entry{Key: clientEntry.Key, Value: clientEntry.Value}, LeaderId: *leaderId, PrevLogIndex: s.lastApplied - 1, PrevLogTerm: s.log[s.lastApplied-1].Term, LeaderCommit: s.commitIndex, Term: s.currentTerm})
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
	for i := 0; i < len(s.peerPorts); i++ {
		resp := <-respChan
		if resp.Succeeds {
			successCounter++
		}
		if successCounter == len(s.peerPorts)/2 {
			//+1 is yourself
			s.commitIndex = s.lastApplied
			//send heartbeat in go routine -> TODO
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
	//check leader's term
	if in.Term < s.currentTerm {
		return &rt.ResultAppend{Succeeds: false, Term: s.currentTerm}, fmt.Errorf("i don't think you are a leader")
	}
	//check if log is consistent
	fmt.Println(s.log, in.PrevLogTerm, "Printing indexes")
	if s.log[in.PrevLogIndex].Term != in.PrevLogTerm {
		return &rt.ResultAppend{Succeeds: false, Term: s.currentTerm}, fmt.Errorf("log is not consistent")
	}
	if in.Entry != nil {
		s.log = append(s.log, Entry{Key: in.Entry.Key, Value: in.Entry.Value, Term: s.currentTerm})
		fmt.Println("I am a follower, I have appended to my log", in.Entry.Key, in.Entry.Value, "with term", s.currentTerm, "my logs are", s.log)
		return &rt.ResultAppend{Succeeds: true, Term: s.currentTerm}, nil
	}
	//is current term is less or equal to leader's term
	s.currentTerm = in.Term
	//Update lastApplied
	if s.commitIndex < in.LeaderCommit {
		s.lastApplied = int64(len(s.log) - 1)
		s.commitIndex = int64(math.Min(float64(in.LeaderCommit), float64(len(s.log))))
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

	server := &RaftServer{leaderId: *leaderId, currentTerm: 0, commitIndex: 0, lastApplied: 0, log: []Entry{{Key: "0", Value: 0, Term: 0}}} //log[0] is dummy
	addresses := strings.Split(*peerPorts, ",")
	for _, addr := range addresses {
		port, _ := strconv.Atoi(addr)
		server.peerPorts = append(server.peerPorts, int64(port))
		//server.nextIndex = append(server.nextIndex, 1)
	}
	// go func() {
	// 	for {
	// 		if *portR == *leaderId {
	// 			for i := 0; i < len(addresses); i++ {
	// 				fmt.Println("Checking my indexLog[]", server.nextIndex)
	// 				server.nextIndex[i] = int64(len(server.log))
	// 			}
	// 		} else {
	// 			server.state = "follower"
	// 		}
	// 		time.Sleep(3 * time.Second)
	// 	}
	// }()

	//temporary
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
