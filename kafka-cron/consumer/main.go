package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"kafka-cron/utils"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-launcher-go/launcher"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

const (
	DefaultKafkaTopic    = "cluster-a-topic"
	DefaultConsumerGroup = "cluster-a"
)

var (
	consumerGroup = flag.String("group", DefaultConsumerGroup, "The name of the consumer group, used for coordination and load balancing")
	kafkaTopic    = flag.String("topic", DefaultKafkaTopic, "topic to consume")
	kafkaBroker   = flag.String("broker", "localhost:9092", "broker in the Kafka cluster")
	pollDuration  = flag.Int("poll", 3000, "The duration of the polling timeout")
)

func main() {
	// enable multi-span attributes
	err := godotenv.Load()
	if err != nil {
		os.Stdout.WriteString("Warning: No .env file found. Consider creating one\n")
	}

	apikey, apikeyPresent := os.LookupEnv("HONEYCOMB_API_KEY")

	if apikeyPresent {
		serviceName, _ := os.LookupEnv("OTEL_SERVICE_NAME")
		os.Stderr.WriteString(fmt.Sprintf("Sending to Honeycomb with API Key <%s> and service name %s\n", apikey, serviceName))

		otelShutdown, err := launcher.ConfigureOpenTelemetry(
			honeycomb.WithApiKey(apikey),
			launcher.WithServiceName(serviceName),
		)
		if err != nil {
			log.Fatalf("error setting up OTel SDK - %e", err)
		}
		defer otelShutdown()
	} else {
		os.Stdout.WriteString("Honeycomb API key not set - disabling OpenTelemetry")
	}

	tracer := otel.Tracer("")
	setupPrometheus(2112)
	flag.Parse()

	// Create producer for retries
	retryProducer, err := utils.SetupProducer(*kafkaBroker)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Delivery confirmation topic_retries
	go func() {
		DeliveryToKafka(retryProducer)
	}()
	// Create consumer
	consumer := SetupConsumer()
	// Subscribe to the topic
	if err := consumer.Subscribe(*kafkaTopic, nil); err != nil {
		fmt.Printf("There was an error subscribing to the topic :\n\t%v\n", err)
		os.Exit(1)
	}
	// Start consuming
	run := true
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(*pollDuration)
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
				traceID, err := trace.TraceIDFromHex(cronJob.TraceID)
				if err != nil {
					fmt.Println(err)
				}
				spanCtx := trace.SpanContextFromContext(context.Background()).WithTraceID(traceID)
				fmt.Println(cronJob.TraceID, "TRACE ID CONSUMER CRONJOB")
				ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)
				parentCtx, parentSpan := tracer.Start(ctx, "consumer_job_received")
				fmt.Println(parentSpan.SpanContext().TraceID(), "SPAN TRACE ID CONSUMER")
				out, err := ExecJob(parentCtx, traceID, tracer, cronJob.Command, cronJob.Args)

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
							_, span := tracer.Start(parentCtx, "produce_job_retry")
							err = retryProducer.Produce(&message, nil)
							if err != nil {
								//Prometheus
								CounterMessagesError.WithLabelValues(*km.TopicPartition.Topic, "message_produced_to_retries_topic_error").Inc()
								fmt.Printf("Failed to produce message: %s\n", err.Error())
							}
							CounterMessagesSuccess.WithLabelValues(*km.TopicPartition.Topic, "message_produced_to_retries_topic").Inc()
							fmt.Println("🤞 Retrying job", cronJob.Retries, "retries left")
							span.End()
						}()
					} else {
						fmt.Println("No retries left")
						//Prometheus
						CounterOfExceededRetries.WithLabelValues(*km.TopicPartition.Topic).Inc()
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
				parentSpan.End()
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
	consumer.Close()
}
