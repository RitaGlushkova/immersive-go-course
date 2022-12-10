package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type cronjob struct {
	Crontab string   `json:"crontab"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	// cluster string
	// retries int
}

const (
	DefaultKafkaTopic    = "test_topic"
	DefaultConsumerGroup = "myGroup"
)

var (
	consumerGroup = flag.String("group", DefaultConsumerGroup, "The name of the consumer group, used for coordination and load balancing")
	kafkaTopic    = flag.String("topic", DefaultKafkaTopic, "The comma-separated list of topics to consume")
)

func main() {
	flag.Parse()

	// Store the config
	cm := kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"group.id":           *consumerGroup,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest",
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

			ev := c.Poll(1000)
			if ev == nil {
				// the Poll timed out and we got nothing'
				fmt.Printf("……\n")
				continue
			}
			// The poll pulled an event, let's now
			// look at the type of Event we've received
			switch e := ev.(type) {

			case *kafka.Message:
				// It's a message
				km := ev.(*kafka.Message)
				cronJob := cronjob{}
				fmt.Printf("✅ Message '%v' received from topic '%v' (partition %d at offset %d) key %v\n",
					string(km.Value),
					string(*km.TopicPartition.Topic),
					km.TopicPartition.Partition,
					km.TopicPartition.Offset,
					string(km.Key))
				err := json.Unmarshal(km.Value, &cronJob)
				if err != nil {
					fmt.Println(err)
				}
				out, err := execJob(cronJob.Command, cronJob.Args)
				if err != nil {
					fmt.Println(err)
				}
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
	fmt.Println("Command Successfully Executed")
	return stdout, nil
}
