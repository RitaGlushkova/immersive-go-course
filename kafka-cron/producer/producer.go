package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
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

	// Store the config
	c := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"acks":              "all"}

	// Create the producer. Variable p holds the new Producer instance.
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

	// Create topic if needed
	CreateTopic(p, topic)

	//Handle eny events that come back from the producer
	go func() {
		//true
		for e := range p.Events() {
			// The `select` blocks until one of the `case` conditions
			// are met - therefore we run it in a Go Routine.
			switch ev := e.(type) {
			case *kafka.Message:
				// It's a delivery report
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
	cronjobs, err := readCrontabfile("crontab.txt")
	if err != nil {
		log.Fatal(err)
	}
	for _, job := range cronjobs {
		myJob := job
		_, er := cron.AddFunc(job.Crontab, func() {
			recordValue, _ := json.Marshal(&myJob)
			message := kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Key:            []byte(uuid.New().String()),
				Value:          []byte(recordValue),
			}
			if err = p.Produce(&message, nil); err != nil {
				fmt.Printf("Failed to produce message: %s\n", err.Error())
			}
		})
		if er != nil {
			fmt.Println(er)
		}
		fmt.Printf("cronjobs: started cron for %+v\n", myJob)
	}
	cron.Start()
	time.Sleep(1 * time.Minute)
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
func CreateTopic(p *kafka.Producer, topic string) {

	a, err := kafka.NewAdminClientFromProducer(p)
	if err != nil {
		fmt.Printf("Failed to create new admin client from producer: %s", err)
		os.Exit(1)
	}
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create topics on cluster.
	// Set Admin options to wait up to 60s for the operation to finish on the remote cluster
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		fmt.Printf("ParseDuration(60s): %s", err)
		os.Exit(1)
	}
	results, err := a.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     2,
			ReplicationFactor: 1}},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		fmt.Printf("Admin Client request error: %v\n", err)
		os.Exit(1)
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			fmt.Printf("Failed to create topic: %v\n", result.Error)
			os.Exit(1)
		}
		fmt.Printf("%v\n", result)
	}
	a.Close()

}
