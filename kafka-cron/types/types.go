package types

import (
	"time"
)

type Cronjob struct {
	Crontab            string      `json:"crontab"`
	Command            string      `json:"command"`
	Args               []string    `json:"args"`
	Cluster            string      `json:"cluster"`
	Retries            int         `json:"retries"`
	TimestampProduced  time.Time   `json:"timestamp"`
	TimestampAttempted []time.Time `json:"timestamp_attempted"`
	TraceID            string      `json:"trace_id"`
}

// type Cronjob struct {
// 	Crontab            string      `json:"crontab"`
// 	Command            string      `json:"command"`
// 	Args               []string    `json:"args"`
// 	Cluster            string      `json:"cluster"`
// 	Retries            int         `json:"retries"`
// 	TimestampProduced  time.Time   `json:"timestamp"`
// 	TimestampAttempted []time.Time `json:"timestamp_attempted"`
// }

type TopicConfig struct {
	TopicNames        []string `json:"topic_name"`
	TopicPartitions   []int    `json:"topic_partitions"`
	TopicReplications []int    `json:"topic_replication_factor"`
}
