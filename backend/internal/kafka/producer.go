package kafka

import (
	"fmt"

	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(brokers, topic string) (*KafkaProducer, error) {
	// Create a new producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		// bootstrap.servers should be localhost:<PLAINTEXT PORTS>
		"bootstrap.servers": brokers,
	})

	// Check for errors
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
		return nil, err
	}
	return &KafkaProducer{producer: p, topic: topic}, nil
}

// Produce sends a message to the Kafka topic
func (kp *KafkaProducer) Produce(message string) error {
	// Create a delivery channel for the producer
	deliveryChan := make(chan kafka.Event)

	// Produce the message to Kafka
	err := kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)
	if err != nil {
		return err
	}

	// Wait for the delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)

	// Check for any errors in the delivery
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}
	close(deliveryChan)
	return nil
}

// Close closes the Kafka producer
func (kp *KafkaProducer) Close() {
	kp.producer.Close()
}
