package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	MessageCounterError = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "producer_error_counter",
		Help: "metric that counts errors in the producer",
	}, []string{
		"topic", "error_type",
	})

	MessageCounterSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "producer_success_counter",
		Help: "metric that counts successes in the producer",
	}, []string{
		"topic", "job_type",
	})
	CronErrorCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cron_error_counter",
		Help: "metric that counts errors in the cron",
	}, []string{
		"cluster", "job_type",
	})
	CronSuccessCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cron_success_counter",
		Help: "metric that counts successes in the cron",
	}, []string{
		"cluster", "job_type",
	})
	LatencyMessageProduced = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "producer_message_latency_produced",
		Help:    "metric that tracks the latency of producing messages",
		Buckets: prometheus.DefBuckets,
	}, []string{
		"topic", "error_type",
	})
	LatencyMessageDelivered = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "producer_message_latency_delivered",
		Help:    "metric that tracks the latency of delivering messages",
		Buckets: prometheus.DefBuckets,
	}, []string{
		"topic", "error_type",
	})
)

func setupPrometheus(port int) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
}
