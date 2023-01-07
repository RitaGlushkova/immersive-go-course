package utils

import (
	"fmt"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func SetupProducer(kafkaBroker string) (*kafka.Producer, error) {
	c := kafka.ConfigMap{
		"bootstrap.servers":   kafkaBroker,
		"delivery.timeout.ms": 10000,
		"acks":                "all"}
	p, errP := kafka.NewProducer(&c)
	if errP != nil {
		if ke, ok := errP.(kafka.Error); ok {
			switch ec := ke.Code(); ec {
			case kafka.ErrInvalidArg:
				return nil, fmt.Errorf("can't create the producer due to configurations: %v", errP)

			default:
				return nil, fmt.Errorf("can't create the producer: %v", errP)

			}
		} else {
			return nil, fmt.Errorf("there's a generic error creating the Producer %v", errP)

		}
	}
	return p, nil
}
