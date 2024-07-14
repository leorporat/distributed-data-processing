// internal/kafka/consumer.go

package kafka

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topics   []string
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(brokers string, groupID string, topics []string) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		// bootstrap.servers should be localhost:<PLAINTEXT PORTS>
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to produce Consumer %s", err)
		os.Exit(1)
		return nil, err
	}

	return &KafkaConsumer{consumer: c, topics: topics}, nil
}

// Consume starts consuming messages from the Kafka topics
func (kc *KafkaConsumer) Consume(messageChan chan<- string) error {

	err := kc.consumer.SubscribeTopics(kc.topics, nil)
	if err != nil {
		return err
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v, terminating program", sig)
			run = false
		default:
			msg, err := kc.consumer.ReadMessage(-1)
			if err != nil {
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
				continue
			}
			messageChan <- string(msg.Value)
		}
	}

	return nil
}

// Close closes the Kafka consumer
func (kc *KafkaConsumer) Close() {
	kc.consumer.Close()
}
