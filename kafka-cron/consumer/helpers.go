package main

import (
	"encoding/json"
	"fmt"
	"kafka-cron/types"
	"kafka-cron/utils"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func ExecJob(command string, args []string) ([]byte, error) {
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
