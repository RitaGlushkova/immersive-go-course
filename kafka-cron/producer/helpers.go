package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/shlex"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func ReadCrontabfile(path string) ([]cronjob, error) {
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
