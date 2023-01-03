package main

import (
	"encoding/json"
	"flag"
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
	Cluster string   `json:"cluster"`
	Retries int      `json:"retries"`
}

var (
	kafkaBroker = flag.String("broker", "localhost:9092", "The comma-separated list of brokers in the Kafka cluster")
)

func main() {
	setupPrometheus(2112)
	flag.Parse()
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
	cronjobs, err := ReadCrontabfile("/cronfile.txt")
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
