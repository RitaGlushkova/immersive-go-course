package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	CounterMessagesError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "message_counter_error",
		Help: "metric that tracks the errors in the consumer or producer for retrying",
	}, []string{
		"topic", "error_type",
	})
	CounterMessagesSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "message_counter_success",
		Help: "metric that tracks the success in the consumer or producer for retrying",
	}, []string{
		"topic", "job_type",
	})
	LatencyExecution = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "consumer_execution_latency",
		Help:    "metric that tracks the latency of executing jobs",
		Buckets: prometheus.DefBuckets,
	}, []string{
		"topic",
	})
	LatencyExecutionSuccess = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "consumer_execution_latency_success",
		Help:    "metric that tracks the latency of successfully executing jobs",
		Buckets: prometheus.DefBuckets,
	}, []string{
		"topic",
	})
	LatencyExecutionError = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "consumer_execution_latency_error",
		Help:    "metric that tracks the latency of failed executing jobs",
		Buckets: prometheus.DefBuckets,
	}, []string{
		"topic",
	})
	MessagesInFlight = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "messages_in_flight",
		Help: "metric that tracks the number of messages in flight",
	}, []string{
		"topic",
	})
	CounterOfExceededRetries = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "counter_of_exceeded_retries",
		Help: "metric that tracks the number of messages that exceeded the number of retries",
	}, []string{
		"topic",
	})
)

func setupPrometheus(port int) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
}
