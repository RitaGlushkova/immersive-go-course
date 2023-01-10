package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	pb "github.com/RitaGlushkova/immersive-go-course/grpc-client-server/prober"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	port         = flag.Int("port", 50051, "The server port")
	LatencyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "golang",
			Name:      "latency_gauge",
			Help:      "metric that tracks the latency",
		}, []string{
			"endpoint",
		})
)

type server struct {
	pb.UnimplementedProberServer
}

func (s *server) DoProbes(ctx context.Context, in *pb.ProbeRequest) (*pb.ProbeReply, error) {
	var sumOfelapsedMsecs = float32(0)
	numberOfRepeats := in.GetNumberOfRequestsToMake()
	var replies = make([]*pb.Reply, 0)
	for i := 0; i < int(numberOfRepeats); i++ {
		var reply pb.Reply
		start := time.Now()
		res, err := http.Get(in.GetEndpoint())
		if err != nil {
			reply.ErrorMessage = fmt.Sprintf("Error: %v, request number %d", err, i)
		} else {
			reply.ReplyCode = int64(res.StatusCode)
		}
		elapsed := time.Since(start)
		elapsedMsecs := float32(elapsed) / float32(time.Millisecond)
		reply.LatencyMsecs = elapsedMsecs

		LatencyGauge.WithLabelValues(in.GetEndpoint()).Set(float64(elapsedMsecs))
		sumOfelapsedMsecs += elapsedMsecs
		replies = append(replies, &reply)
	}
	averageLatencyMsecs := sumOfelapsedMsecs / float32(numberOfRepeats)
	return &pb.ProbeReply{AverageLatencyMsecs: averageLatencyMsecs, Replies: replies}, nil
}

func init() {
	// Metrics have to be registered to be exposed:
	//It only can be run once !!
	prometheus.MustRegister(LatencyGauge)
	http.Handle("/metrics", promhttp.Handler())
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
