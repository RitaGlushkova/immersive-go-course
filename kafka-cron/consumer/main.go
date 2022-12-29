package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type cronjob struct {
	Crontab string   `json:"crontab"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Cluster string   `json:"cluster"`
	Retries int      `json:"retries"`
}

const (
	DefaultKafkaTopic    = "test_topic"
	DefaultConsumerGroup = "myGroup"
)

var (
	consumerGroup = flag.String("group", DefaultConsumerGroup, "The name of the consumer group, used for coordination and load balancing")
	kafkaTopic    = flag.String("topic", DefaultKafkaTopic, "The comma-separated list of topics to consume")
	kafkaBroker   = flag.String("broker", "localhost:9092", "The comma-separated list of brokers in the Kafka cluster")
	ErrorCounter  = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "consumer_error_counter",
		Help: "metric that tracks the error",
	}, []string{
		"topic", "error_type",
	})
	SuccessCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "consumer_success_counter",
		Help: "metric that tracks the success in the consumer",
	}, []string{
		"topic", "type",
	})
	LatencyExecution = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "consumer_execution_latency",
		Help:    "metric that tracks the latency of executing commands",
		Buckets: prometheus.DefBuckets,
	}, []string{
		"topic",
	})
)

func init() {
	http.Handle("/metrics", promhttp.Handler())
}

func main() {
	flag.Parse()
	_, err := setupPrometheus(2112)
	if err != nil {
		log.Fatal("Failed to listen on port :2112", err)
	}
	p, errP := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": *kafkaBroker})
	// Check for errors
	if errP != nil {
		if ke, ok := errP.(kafka.Error); ok {
			switch ec := ke.Code(); ec {
			case kafka.ErrInvalidArg:
				fmt.Printf("Can't create the producer because you've configured it wrong (code: %d)!\n\t%v\n", ec, errP)
				os.Exit(1)
			default:
				fmt.Printf("Can't create the producer (code: %d)!\n\t%v\n", ec, errP)
				os.Exit(1)
			}
		} else {
			fmt.Printf("There's a generic error creating the Producer! %v", errP.Error())
			os.Exit(1)
		}
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:

				if ev.TopicPartition.Error != nil {
					//Prometheus
					ErrorCounter.WithLabelValues(*ev.TopicPartition.Topic, "producer_retries_delivery").Inc()
					fmt.Printf("Failed to send message '%v' to topic '%v'\n\tErr: %v",
						string(ev.Value),
						string(*ev.TopicPartition.Topic),
						ev.TopicPartition.Error)
				} else {
					//Prometheus
					SuccessCounter.WithLabelValues(*ev.TopicPartition.Topic, "producer_retries_delivery").Inc()
					fmt.Printf("✅ Message '%v' with key '%v' delivered to topic '%v' (partition %d at offset %d)\n",
						string(ev.Value),
						string(ev.Key),
						string(*ev.TopicPartition.Topic),
						ev.TopicPartition.Partition,
						ev.TopicPartition.Offset)
					fmt.Println(ev.TopicPartition)
				}
			case kafka.Error:
				ErrorCounter.WithLabelValues(*kafkaTopic, "producer_retries_delivery").Inc()
				// It's an error
				fmt.Printf("Caught an error:\n\t%v\n", ev.Error())
			default:
				ErrorCounter.WithLabelValues(*kafkaTopic, "producer_retries_delivery").Inc()
				// It's not anything we were expecting
				fmt.Printf("Got an event that's not a Message or Error 👻\n\t%v\n", ev)

			}
		}
	}()

	// Store the config
	cm := kafka.ConfigMap{
		"bootstrap.servers":  *kafkaBroker,
		"group.id":           *consumerGroup,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&cm)
	// Check for errors in creating the Consumer
	if err != nil {
		if ke, ok := err.(kafka.Error); ok {
			switch ec := ke.Code(); ec {
			case kafka.ErrInvalidArg:
				fmt.Printf("Can't create the Consumer because you've configured it wrong (code: %d)!\n\t%v\n", ec, err)
				os.Exit(1)
			default:
				fmt.Printf("Can't create the Consumer (Kafka error code %d)\n\tError: %v\n", ec, err)
				os.Exit(1)
			}
		} else {
			// It's not a kafka.Error
			fmt.Printf("There's a generic error creating the Consumer! %v", err.Error())
			os.Exit(1)
		}

	}
	fmt.Printf("Created Consumer %v\n", cm)

	// Subscribe to the topic
	if err := c.Subscribe(*kafkaTopic, nil); err != nil {
		fmt.Printf("There was an error subscribing to the topic :\n\t%v\n", err)
		os.Exit(1)
	}

	run := true

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:

			ev := c.Poll(3000)
			if ev == nil {
				// the Poll timed out and we got nothing'
				fmt.Printf("……\n")
				continue
			}
			// The poll pulled an event, let's now
			// look at the type of Event we've received
			switch e := ev.(type) {

			case *kafka.Message:
				//Prometheus
				start := time.Now()
				// It's a message
				km := ev.(*kafka.Message)
				cronJob := cronjob{}
				fmt.Printf("💌 Message '%v' received from topic '%v' (partition %d at offset %d) key %v\n",
					string(km.Value),
					string(*km.TopicPartition.Topic),
					km.TopicPartition.Partition,
					km.TopicPartition.Offset,
					string(km.Key))
				SuccessCounter.WithLabelValues(*km.TopicPartition.Topic, "consumer_success_message_received").Inc()
				err := json.Unmarshal(km.Value, &cronJob)
				if err != nil {
					//Prometheus
					ErrorCounter.WithLabelValues(*km.TopicPartition.Topic, "consumer_unmarshal").Inc()
					fmt.Println(err)
				}
				out, err := execJob(cronJob.Command, cronJob.Args)
				LatencyExecution.WithLabelValues(*km.TopicPartition.Topic).Observe(time.Since(start).Seconds())
				if err != nil {
					//Prometheus
					ErrorCounter.WithLabelValues(*km.TopicPartition.Topic, "consumer_execution").Inc()
					fmt.Println("😿 Error executing job", err)
					fmt.Println(cronJob.Retries, "retries left")
					if cronJob.Retries > 0 {
						cronJob.Retries = cronJob.Retries - 1
						fmt.Println(cronJob.Retries, "Reties before sending to kafka")
						recordValue, _ := json.Marshal(&cronJob)
						var topic string
						if strings.Contains(*kafkaTopic, "_retries") {
							topic = *kafkaTopic
						} else {
							topic = fmt.Sprintf("%v_retries", *kafkaTopic)
						}
						message := kafka.Message{
							TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
							Key:            []byte(km.Key),
							Value:          []byte(recordValue),
						}
						err = p.Produce(&message, nil)
						if err != nil {
							//Prometheus
							ErrorCounter.WithLabelValues(*km.TopicPartition.Topic, "produce_retries_produce_message").Inc()
							fmt.Printf("Failed to produce message: %s\n", err.Error())
						}
						fmt.Println("🤞 Retrying job", cronJob.Retries, "retries left")
						time.Sleep(5 * time.Second)
					} else {
						fmt.Println("No retries left")
					}
				}
				SuccessCounter.WithLabelValues(*km.TopicPartition.Topic, "consumer_success_job_executed").Inc()
				fmt.Println(string(out))
			case kafka.Error:
				// It's an error
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				// It's not anything we were expecting
				fmt.Printf("Ignored event \n\t%v\n", ev)

			}
		}
	}
	fmt.Printf("👋 … and we're done. Closing the consumer and exiting.\n")

	// Now we can exit
	c.Close()
}

func execJob(command string, args []string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println("😻 Command Successfully Executed")
	return stdout, nil
}

func setupPrometheus(port int) (int, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, err
	}
	go http.Serve(lis, nil)
	return lis.Addr().(*net.TCPAddr).Port, nil
}
