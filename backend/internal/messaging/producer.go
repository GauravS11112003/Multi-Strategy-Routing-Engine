package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

var Producer sarama.SyncProducer

func ConnectProducer(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = 100 * time.Millisecond

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("create kafka producer: %w", err)
	}

	Producer = producer
	log.Println("Connected to Kafka (producer)")
	return nil
}

func CloseProducer() {
	if Producer != nil {
		Producer.Close()
		log.Println("Kafka producer closed")
	}
}

func PublishEvent(topic string, eventType EventType, data interface{}) error {
	if Producer == nil {
		return nil
	}

	event := Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.ID),
		Value: sarama.ByteEncoder(payload),
	}

	partition, offset, err := Producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	log.Printf("Published %s to %s [partition=%d, offset=%d]", eventType, topic, partition, offset)
	return nil
}
