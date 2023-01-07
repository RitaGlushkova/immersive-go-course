package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"kafka-cron/utils"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

const (
	DefaultKafkaTopic    = "cluster-a-topic"
	DefaultConsumerGroup = "cluster-a"
)

var (
	consumerGroup = flag.String("group", DefaultConsumerGroup, "The name of the consumer group, used for coordination and load balancing")
	kafkaTopic    = flag.String("topic", DefaultKafkaTopic, "The comma-separated list of topics to consume")
	kafkaBroker   = flag.String("broker", "localhost:9092", "The comma-separated list of brokers in the Kafka cluster")
)

func main() {
	setupPrometheus(2112)
	flag.Parse()

	// Create producer for retries
	p, err := utils.SetupProducer(*kafkaBroker)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Produce messages to topic_retries
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:

				if ev.TopicPartition.Error != nil {
					CounterMessagesError.WithLabelValues(*ev.TopicPartition.Topic, "retries_topic_partition_error").Inc()
					utils.PrintDeliveryFairure(ev)
				} else {
					CounterMessagesSuccess.WithLabelValues(*ev.TopicPartition.Topic, "message_delivered_to_retries_topic").Inc()
					utils.PrintDeliveryConfirmation(ev)
				}
			case kafka.Error:
				CounterMessagesError.WithLabelValues(*kafkaTopic, "producer_events_retries_kafka_error").Inc()
				fmt.Printf("Caught an error:\n\t%v\n", ev.Error())
			default:
				CounterMessagesError.WithLabelValues(*kafkaTopic, "producer_events_retries_default_error").Inc()
				// It's not anything we were expecting
				fmt.Printf("Got an event that's not a Message or Error 👻\n\t%v\n", ev)

			}
		}
	}()

	// Configure Consumer
	cm := kafka.ConfigMap{
		"bootstrap.servers":  *kafkaBroker,
		"group.id":           *consumerGroup,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "latest",
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
				km := ev.(*kafka.Message)
				cronJob := ReceiveMessage(km)
				// Prometheus
				if strings.Contains(*km.TopicPartition.Topic, "retries") {
					LatencyAttemptedRetryToReception.WithLabelValues(*km.TopicPartition.Topic).Observe(time.Since(cronJob.TimestampAttempted[len(cronJob.TimestampAttempted)-1]).Seconds())
				} else {
					LatencyProductionToReception.WithLabelValues(*km.TopicPartition.Topic).Observe(time.Since(cronJob.TimestampProduced).Seconds())
				}
				startExec := time.Now()
				out, err := ExecJob(cronJob.Command, cronJob.Args)
				// Prometheus
				LatencyExecution.WithLabelValues(*km.TopicPartition.Topic, cronJob.Command).Observe(time.Since(startExec).Seconds())
				if err != nil {
					//Prometheus
					CounterMessagesError.WithLabelValues(*km.TopicPartition.Topic, "consumer_job_execution_error").Inc()
					fmt.Println("😿 Error executing job", err)
					if cronJob.Retries > 0 {
						go func() {
							time.Sleep(1 * time.Second)
							cronJob.Retries = cronJob.Retries - 1
							fmt.Println(cronJob.Retries, "retries left")
							cronJob.TimestampAttempted = append(cronJob.TimestampAttempted, time.Now())
							recordValue, _ := json.Marshal(&cronJob)
							retryTopic := CreateRetryTopic()
							message := kafka.Message{
								TopicPartition: kafka.TopicPartition{Topic: &retryTopic, Partition: kafka.PartitionAny},
								Key:            []byte(km.Key),
								Value:          []byte(recordValue),
							}
							err = p.Produce(&message, nil)
							if err != nil {
								//Prometheus
								CounterMessagesError.WithLabelValues(*km.TopicPartition.Topic, "retries_topic_Produce_message_error").Inc()
								fmt.Printf("Failed to produce message: %s\n", err.Error())
							}
							fmt.Println("🤞 Retrying job", cronJob.Retries, "retries left")
						}()
					} else {
						fmt.Println("No retries left")
						//Prometheus
						CounterOfExceededRetries.WithLabelValues(*km.TopicPartition.Topic, cronJob.Command).Inc()
					}
					//Prometheus
					LatencyExecutionError.WithLabelValues(*km.TopicPartition.Topic, cronJob.Command).Observe(time.Since(startExec).Seconds())
				} else {
					//Prometheus
					CounterMessagesSuccess.WithLabelValues(*km.TopicPartition.Topic, "consumer_success_job_executed").Inc()
					//Print successful output of the job
					fmt.Println(string(out))
					//Prometheus
					LatencyExecutionSuccess.WithLabelValues(*km.TopicPartition.Topic, cronJob.Command).Observe(time.Since(startExec).Seconds())
				}
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
				//Prometheus
				CounterMessagesError.WithLabelValues(*kafkaTopic, "consumer_kafka_error").Inc()
			default:
				// It's not anything we were expecting
				fmt.Printf("Ignored event \n\t%v\n", ev)
				//Prometheus
				CounterMessagesError.WithLabelValues(*kafkaTopic, "consumer_ignored_event").Inc()
			}
		}
	}
	fmt.Printf("👋 … and we're done. Closing the consumer and exiting.\n")
	c.Close()
}
