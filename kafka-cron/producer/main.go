package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/shlex"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type cronjob struct {
	Crontab string   `json:"crontab"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Cluster string   `json:"cluster"`
	Retries int      `json:"retries"`
}

var (
	kafkaBroker  = flag.String("broker", "localhost:9092", "The comma-separated list of brokers in the Kafka cluster")
	LatencyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "golang",
			Name:      "latency_gauge",
			Help:      "metric that tracks the latency",
		}, []string{
			"endpoint",
		})
)

func init() {
	// Metrics have to be registered to be exposed:
	//It only can be run once !!
	prometheus.MustRegister(LatencyGauge)
	http.Handle("/metrics", promhttp.Handler())
}

func main() {
	flag.Parse()
	_, err := setupPrometheus(2112)
	if err != nil {
		log.Fatal("Failed to listen on port :2112", err)
	}
	topics := []string{"test_topic", "my_topic", "test_topic_retries", "my_topic_retries"}
	partitions := []int{1, 2, 1, 1}
	replicas := []int{1, 1, 1, 1}
	// Store the config
	c := kafka.ConfigMap{
		"bootstrap.servers":   *kafkaBroker,
		"delivery.timeout.ms": 10000,
		"acks":                "all"}

	// Create the producer
	p, err := kafka.NewProducer(&c)

	// Check for errors
	if err != nil {
		if ke, ok := err.(kafka.Error); ok {
			switch ec := ke.Code(); ec {
			case kafka.ErrInvalidArg:
				fmt.Printf("Can't create the producer because you've configured it wrong (code: %d)!\n\t%v\n", ec, err)
				os.Exit(1)
			default:
				fmt.Printf("Can't create the producer (code: %d)!\n\t%v\n", ec, err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("There's a generic error creating the Producer! %v", err.Error())
			os.Exit(1)
		}

	}

	// Create topic
	err = CreateTopic(p, topics, partitions, replicas)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:

				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to send message '%v' to topic '%v'\n\tErr: %v",
						string(ev.Value),
						string(*ev.TopicPartition.Topic),
						ev.TopicPartition.Error)
				} else {
					fmt.Printf("✅ Message '%v' with key '%v' delivered to topic '%v' (partition %d at offset %d)\n",
						string(ev.Value),
						string(ev.Key),
						string(*ev.TopicPartition.Topic),
						ev.TopicPartition.Partition,
						ev.TopicPartition.Offset)
					fmt.Println(ev.TopicPartition)
				}
			case kafka.Error:
				// It's an error
				fmt.Printf("Caught an error:\n\t%v\n", ev.Error())
			default:
				// It's not anything we were expecting
				fmt.Printf("Got an event that's not a Message or Error 👻\n\t%v\n", ev)

			}
		}
	}()

	log.Info("Create new cron")
	cron := cron.New(cron.WithSeconds())
	cronjobs, err := readCrontabfile("/cronfile.txt")
	if err != nil {
		log.Fatal(err)
	}

	for _, job := range cronjobs {
		myJob := job
		_, er := cron.AddFunc(job.Crontab, func() {

			if myJob.Cluster == "cluster-a" {
				recordValue, _ := json.Marshal(&myJob)
				message1 := kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topics[0], Partition: kafka.PartitionAny},
					Key:            []byte(uuid.New().String()),
					Value:          []byte(recordValue),
				}
				err = p.Produce(&message1, nil)
				if err != nil {
					fmt.Printf("Failed to produce message: %s\n", err.Error())
				}
			}
			if myJob.Cluster == "cluster-b" {
				recordValue, _ := json.Marshal(&myJob)
				message2 := kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topics[1], Partition: kafka.PartitionAny},
					Key:            []byte(uuid.New().String()),
					Value:          []byte(recordValue),
				}
				err = p.Produce(&message2, nil)
				if err != nil {
					fmt.Printf("Failed to produce message: %s\n", err.Error())
				}
			}
		})
		if er != nil {
			fmt.Println(er)
		}
		fmt.Printf("cronjobs: started cron for %+v\n", myJob)
	}
	cron.Run()
	//time.Sleep(1 * time.Minute)
	fmt.Printf("Flushing outstanding messages\n")
	// Flush the Producer queue
	t := 10000
	if r := p.Flush(t); r > 0 {
		fmt.Printf("\n--\n⚠️ Failed to flush all messages after %d milliseconds. %d message(s) remain\n", t, r)
	} else {
		fmt.Println("\n--\n✨ All messages flushed from the queue")
	}

	// Now we can exit
	p.Close()
}

// CreateTopic creates a topic using the Admin Client API
func CreateTopic(p *kafka.Producer, topics []string, partitions, replicas []int) error {

	a, err := kafka.NewAdminClientFromProducer(p)
	if err != nil {
		return fmt.Errorf("failed to create new admin client from producer: %s", err)

	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create topics on cluster.
	// Set Admin options to wait up to 60s for the operation to finish on the remote cluster
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		return fmt.Errorf("ParseDuration(60s): %s", err)

	}
	var topicsSpec []kafka.TopicSpecification
	for i, topic := range topics {
		var topicSpec = kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     partitions[i],
			ReplicationFactor: replicas[i]}
		topicsSpec = append(topicsSpec, topicSpec)
	}
	results, err := a.CreateTopics(
		ctx,
		topicsSpec,
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		return fmt.Errorf("admin Client request error: %v", err)

	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			return fmt.Errorf("failed to create topic: %v", result.Error)

		}
		fmt.Printf("%v\n", result)
	}
	a.Close()
	return nil
}

func readCrontabfile(path string) ([]cronjob, error) {
	readFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer readFile.Close()
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string
	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	result := make([]cronjob, 0)
	for _, line := range fileLines {
		val, err := shlex.Split(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing line: %v", err)
		}
		retries, err := strconv.Atoi(val[len(val)-1])
		if err != nil {
			return nil, fmt.Errorf("retries arg couldn't be converted to a number: %v", err)
		}
		cj := cronjob{
			Crontab: strings.Join(val[0:6], " "),
			Command: val[6],
			Args:    val[7 : len(val)-2],
			Cluster: val[len(val)-2],
			Retries: retries,
		}
		result = append(result, cj)
	}
	return result, nil
}

func setupPrometheus(port int) (int, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, err
	}
	go http.Serve(lis, nil)
	return lis.Addr().(*net.TCPAddr).Port, nil
}
