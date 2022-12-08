package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type cronjob struct {
	Crontab string   `json:"crontab"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	// cluster string
	// retries int
}

func main() {

	topic := "test_topic"
	// Create Consumer instance

	// Store the config
	cm := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	}
	// sigchan := make(chan os.Signal, 1)
	// signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// Variable p holds the new Consumer instance.
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
	fmt.Printf("Created Consumer %v\n", c)

	// Subscribe to the topic
	if err := c.Subscribe(topic, nil); err != nil {
		fmt.Printf("There was an error subscribing to the topic :\n\t%v\n", err)
		os.Exit(1)
	}

	run := true

	for run {
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
			fmt.Printf("✅ Message '%v' received from topic '%v' (partition %d at offset %d)\n",
				string(km.Value),
				string(*km.TopicPartition.Topic),
				km.TopicPartition.Partition,
				km.TopicPartition.Offset)
			err := json.Unmarshal(km.Value, &cronJob)
			if err != nil {
				fmt.Println(err)
			}
			execJob(cronJob.Command, cronJob.Args)

		case kafka.Error:
			// It's an error
			em := ev.(kafka.Error)
			fmt.Printf("Caught an error:\n\t%v\n", em)
			if e.Code() == kafka.ErrAllBrokersDown {
				run = false
			}
		default:
			// It's not anything we were expecting
			fmt.Printf("Got an event that's not a Message, Error, or PartitionEOF 👻\n\t%v\n", ev)

		}
	}
	fmt.Printf("👋 … and we're done. Closing the consumer and exiting.\n")

	// Now we can exit
	c.Close()
}

func execJob(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Command Successfully Executed")
	fmt.Println(string(stdout))
}
