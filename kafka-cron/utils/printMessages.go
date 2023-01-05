package utils

import (
	"fmt"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func PrintDeliveryConfirmation(ev *kafka.Message) {
	fmt.Printf("✅ Message '%v' with key '%v' delivered to topic '%v' (partition %d at offset %d)\n",
		string(ev.Value),
		string(ev.Key),
		string(*ev.TopicPartition.Topic),
		ev.TopicPartition.Partition,
		ev.TopicPartition.Offset)
	fmt.Println(ev.TopicPartition)
}

func PrintDeliveryFairure(ev *kafka.Message) {
	fmt.Printf("Failed to send message '%v' to topic '%v'\n\tErr: %v",
		string(ev.Value),
		string(*ev.TopicPartition.Topic),
		ev.TopicPartition.Error)
}

func PrintConfirmatonForReceivedMessage(km *kafka.Message) {
	fmt.Printf("💌 Message '%v' received from topic '%v' (partition %d at offset %d) key %v\n",
		string(km.Value),
		string(*km.TopicPartition.Topic),
		km.TopicPartition.Partition,
		km.TopicPartition.Offset,
		string(km.Key))
}
