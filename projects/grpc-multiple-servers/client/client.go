package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/RitaGlushkova/immersive-go-course/projects/grpc-multiple-servers/userinfo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
	dob  = flag.String("dob", "01/01/1990", "Date of birth")
)

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

func contactServer(addr string) (*pb.UserInfoReply, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserInfoClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendUserInfo(ctx, &pb.UserInfoRequest{Name: *name, Dob: *dob})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func main() {
	flag.Parse()
	r, err := contactServer(*addr)
	if err != nil {
		log.Fatalf("Couldn't send request Info: %v", err)
	}
	if r.GetRedirect() != "" {
		log.Printf("Redirect to: %v", r.GetRedirect())
		r, err = contactServer(r.GetRedirect())
		if err != nil {
			log.Fatalf("Couldn't send request Info: %v", err)
		}
		if r.Notfound != "" {
			log.Fatalf("User not found: %v", r.GetNotfound())
		}
	}
	log.Printf("User Info: %v", UserInfo{Name: r.GetName(), DOB: r.GetDob(), Email: r.GetEmail(), College: r.GetCollege(), Occupation: r.GetOccupation(), Age: r.GetAge()})
}
