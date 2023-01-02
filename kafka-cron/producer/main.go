package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	//"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/shlex"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	kafkaBroker = flag.String("broker", "localhost:9092", "The comma-separated list of brokers in the Kafka cluster")

	// Metrics have to be registered to be exposed:
	MessageCounterError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "producer_error_counter",
		Help: "metric that counts errors in the producer",
	}, []string{
		"topic", "error_type",
	})

	MessageCounterSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "producer_success_counter",
		Help: "metric that counts successes in the producer",
	}, []string{
		"topic", "job_type",
	})
	CronErrorCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cron_error_counter",
		Help: "metric that counts errors in the cron",
	}, []string{
		"cluster", "job_type",
	})
	CronSuccessCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cron_success_counter",
		Help: "metric that counts successes in the cron",
	}, []string{
		"cluster", "job_type",
	})
	LatencyMessageProduced = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "producer_message_latency_produced",
		Help:    "metric that tracks the latency of producing messages",
		Buckets: prometheus.DefBuckets,
	}, []string{
		"topic", "error_type",
	})
	LatencyMessageDelivered = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "producer_message_latency_delivered",
		Help:    "metric that tracks the latency of delivering messages",
		Buckets: prometheus.DefBuckets,
	}, []string{
		"topic", "error_type",
	})
)

// func init() {

// 	//It only can be run once !!
// 	http.Handle("/metrics", promhttp.Handler())
// }

func main() {
	setupPrometheus(2112)
	flag.Parse()
	// _, err := setupPrometheus(2112)
	// if err != nil {
	// 	log.Fatal("Failed to listen on port :2112", err)
	// }
	topics := []string{"cluster-a-topic", "cluster-b-topic", "cluster-a-topic-retries", "cluster-b-topic-retries"}
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
					//Prometheus
					MessageCounterError.WithLabelValues(*ev.TopicPartition.Topic, "message_delivery_error").Inc()
					//Print
					fmt.Printf("Failed to send message '%v' to topic '%v'\n\tErr: %v",
						string(ev.Value),
						string(*ev.TopicPartition.Topic),
						ev.TopicPartition.Error)
				} else {
					//Prometheus
					MessageCounterSuccess.WithLabelValues(*ev.TopicPartition.Topic, "message_delivery_to_topic").Inc()
					fmt.Printf("✅ Message '%v' with key '%v' delivered to topic '%v' (partition %d at offset %d)\n",
						string(ev.Value),
						string(ev.Key),
						string(*ev.TopicPartition.Topic),
						ev.TopicPartition.Partition,
						ev.TopicPartition.Offset)
					fmt.Println(ev.TopicPartition)
				}
			case *kafka.Error:
				// It's an error
				fmt.Printf("Caught an error:\n\t%v\n", ev.Error())
				MessageCounterError.WithLabelValues("all_topics", "message_delivery_kafka_error").Inc()
			default:
				// It's not anything we were expecting
				fmt.Printf("Got an event that's not a Message or Error 👻\n\t%v\n", ev)
				MessageCounterError.WithLabelValues("all_topics", "message_delivery_unknown_error").Inc()
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
		_, cronErr := cron.AddFunc(job.Crontab, func() {
			var message kafka.Message
			if myJob.Cluster == "cluster-a" {
				recordValue, _ := json.Marshal(&myJob)
				message = kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topics[0], Partition: kafka.PartitionAny},
					Key:            []byte(uuid.New().String()),
					Value:          []byte(recordValue),
				}
			}
			if myJob.Cluster == "cluster-b" {
				recordValue, _ := json.Marshal(&myJob)
				message = kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topics[1], Partition: kafka.PartitionAny},
					Key:            []byte(uuid.New().String()),
					Value:          []byte(recordValue),
				}
			}
			//Prometheus
			start := time.Now()
			errStr := ""
			err = p.Produce(&message, nil)
			if err != nil {
				//Prometheus
				MessageCounterError.WithLabelValues(*message.TopicPartition.Topic, "producing_message").Inc()
				fmt.Printf("Failed to produce message: %s\n", err.Error())
				errStr = err.Error()
			}
			LatencyMessageProduced.WithLabelValues(*message.TopicPartition.Topic, errStr).Observe(time.Since(start).Seconds())
			MessageCounterSuccess.WithLabelValues(*message.TopicPartition.Topic, "producing_message").Inc()
		})
		if cronErr != nil {
			fmt.Println(cronErr)
			CronErrorCounter.WithLabelValues(myJob.Cluster, "cronjob_error").Inc()
		}
		fmt.Printf("cronjobs: started cron for %+v\n", myJob)
		CronSuccessCounter.WithLabelValues(myJob.Cluster, "cronjob_success").Inc()
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

//	func setupPrometheus(port int) (int, error) {
//		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
//		if err != nil {
//			return 0, err
//		}
//		go http.Serve(lis, nil)
//		return lis.Addr().(*net.TCPAddr).Port, nil
//	}
func setupPrometheus(port int) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
}
