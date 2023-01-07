package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"kafka-cron/types"
	"kafka-cron/utils"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	kafkaBroker = flag.String("broker", "localhost:9092", "The comma-separated list of brokers in the Kafka cluster")
)

func main() {
	setupPrometheus(2112)
	flag.Parse()
	topicConfig := types.TopicConfig{
		TopicNames:        []string{"cluster-a-topic", "cluster-b-topic", "cluster-a-topic-retries", "cluster-b-topic-retries"},
		TopicPartitions:   []int{1, 2, 1, 1},
		TopicReplications: []int{1, 1, 1, 1},
	}

	// Create the producer
	p, err := utils.SetupProducer(*kafkaBroker)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create topic
	err = CreateTopic(p, topicConfig.TopicNames, topicConfig.TopicPartitions, topicConfig.TopicReplications)
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
					utils.PrintDeliveryFairure(ev)
				} else {
					//Prometheus
					MessageCounterSuccess.WithLabelValues(*ev.TopicPartition.Topic, "message_delivery_to_topic").Inc()
					utils.PrintDeliveryConfirmation(ev)
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
		myJob := types.Cronjob{
			Crontab:           job.Crontab,
			Command:           job.Command,
			Args:              job.Args,
			Cluster:           job.Cluster,
			Retries:           job.Retries,
			TimestampProduced: time.Now()}
		_, cronErr := cron.AddFunc(job.Crontab, func() {
			var message kafka.Message
			//add timnestamp
			myJob.TimestampProduced = time.Now()
			if myJob.Cluster == "cluster-a" {
				recordValue, _ := json.Marshal(&myJob)
				message = kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topicConfig.TopicNames[0], Partition: kafka.PartitionAny},
					Key:            []byte(uuid.New().String()),
					Value:          []byte(recordValue),
				}
			}
			if myJob.Cluster == "cluster-b" {
				recordValue, _ := json.Marshal(&myJob)
				message = kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topicConfig.TopicNames[1], Partition: kafka.PartitionAny},
					Key:            []byte(uuid.New().String()),
					Value:          []byte(recordValue),
				}
			}
			//Prometheus
			startProduce := time.Now()
			errStr := ""
			err = p.Produce(&message, nil)
			if err != nil {
				//Prometheus
				MessageCounterError.WithLabelValues(*message.TopicPartition.Topic, "producing_message").Inc()
				fmt.Printf("Failed to produce message: %s\n", err.Error())
				errStr = err.Error()
			}
			LatencyMessageProduced.WithLabelValues(*message.TopicPartition.Topic, errStr).Observe(time.Since(startProduce).Seconds())
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
	fmt.Printf("Flushing outstanding messages\n")
	t := 10000
	if r := p.Flush(t); r > 0 {
		fmt.Printf("\n--\n⚠️ Failed to flush all messages after %d milliseconds. %d message(s) remain\n", t, r)
	} else {
		fmt.Println("\n--\n✨ All messages flushed from the queue")
	}
	p.Close()
}
