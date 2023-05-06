package main

import (
	"context"
	"flag"
	"fmt"
	"hash/fnv"
	"log"
	"net"

	pb "github.com/RitaGlushkova/immersive-go-course/projects/grpc-multiple-servers/userinfo"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedUserInfoServer
}

type UserInfo struct {
	Name       string
	DOB        string
	Email      string
	College    string
	Occupation string
	Age        int32
	Redirect   string
	Port       int32
	Notfound   string
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func locatePort(name, dob string) int {
	//first reject request with & in name or dob
	//reject id length(dob) <= 255 convert into hex
	hash := hash(dob+"&"+name) % 3
	//	hash := hash(dob+"&"+name) % 3
	if hash == 0 {
		return 50051
	}
	if hash == 1 {
		return 50052
	}
	return 50053
}

func (s *server) SendUserInfo(ctx context.Context, in *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	myPort := locatePort(in.GetName(), in.GetDob())
	if *port != myPort {
		return &pb.UserInfoReply{Redirect: fmt.Sprintf("localhost:%d", myPort)}, nil
	}
	log.Printf("Received: %v, %v", in.GetName(), in.GetDob())
	for _, user := range Users {
		if in.GetName() == user.Name && in.GetDob() == user.DOB {
			return &pb.UserInfoReply{Name: user.Name, Dob: user.DOB, Email: user.Email, College: user.College, Occupation: user.Occupation, Age: user.Age}, nil
		}
	}
	return &pb.UserInfoReply{Notfound: "User not found"}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserInfoServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
