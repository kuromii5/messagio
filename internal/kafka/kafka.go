package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

func NewProducer(brokers []string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to initialize kafka producer: %s", err)
	}

	return producer
}

func NewConsumer(brokers []string) sarama.Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to initialize kafka consumer: %s", err)
	}

	return consumer
}
