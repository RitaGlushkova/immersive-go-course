package main

import (
	"context"
	"encoding/json"
	"fmt"
	"kafka-cron/types"
	"kafka-cron/utils"
	"os"
	"os/exec"
	"strings"

	"go.opentelemetry.io/otel/trace"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func ExecJob(parentCtx context.Context, traceID trace.TraceID, tracer trace.Tracer, command string, args []string) ([]byte, error) {
	_, span := tracer.Start(parentCtx, "consumer_job_execution")
	defer span.End()
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

func ReceiveMessage(km *kafka.Message) types.Cronjob {
	MessagesInFlight.WithLabelValues(*kafkaTopic).Inc()
	defer MessagesInFlight.WithLabelValues(*kafkaTopic).Dec()
	//Prometheus
	cronJob := types.Cronjob{}
	utils.PrintConfirmatonForReceivedMessage(km)
	//Prometheus
	CounterMessagesSuccess.WithLabelValues(*km.TopicPartition.Topic, "consumer_message_received").Inc()
	err := json.Unmarshal(km.Value, &cronJob)
	if err != nil {
		//Prometheus
		CounterMessagesError.WithLabelValues(*km.TopicPartition.Topic, "consumer_unmarshal_error").Inc()
		fmt.Println(err)
	}
	return cronJob
}

func CreateRetryTopic() string {
	var topic string
	if strings.Contains(*kafkaTopic, "-retries") {
		topic = *kafkaTopic
	} else {
		topic = fmt.Sprintf("%v-retries", *kafkaTopic)
	}
	return topic
}

func DeliveryToKafka(retryProducer *kafka.Producer) {
	for e := range retryProducer.Events() {
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
}

func SetupConsumer() (consumer *kafka.Consumer) {
	// Configure Consumer
	cm := kafka.ConfigMap{
		"bootstrap.servers":  *kafkaBroker,
		"group.id":           *consumerGroup,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "latest",
		"enable.auto.commit": false,
	}
	consumer, err := kafka.NewConsumer(&cm)
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
	return consumer
}
