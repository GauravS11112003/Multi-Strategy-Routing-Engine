package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type MessageHandler func(event Event) error

type ConsumerGroup struct {
	group    sarama.ConsumerGroup
	handlers map[string]MessageHandler
}

func NewConsumerGroup(brokers []string, groupID string) (*ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("create consumer group: %w", err)
	}

	log.Printf("Connected to Kafka consumer group: %s", groupID)
	return &ConsumerGroup{
		group:    group,
		handlers: make(map[string]MessageHandler),
	}, nil
}

func (cg *ConsumerGroup) RegisterHandler(topic string, handler MessageHandler) {
	cg.handlers[topic] = handler
}

func (cg *ConsumerGroup) Start(ctx context.Context, topics []string) error {
	consumer := &groupHandler{handlers: cg.handlers}

	for {
		if err := cg.group.Consume(ctx, topics, consumer); err != nil {
			return fmt.Errorf("consumer group error: %w", err)
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (cg *ConsumerGroup) Close() error {
	return cg.group.Close()
}

type groupHandler struct {
	handlers map[string]MessageHandler
}

func (h *groupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *groupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to unmarshal message from %s: %v", msg.Topic, err)
			session.MarkMessage(msg, "")
			continue
		}

		handler, ok := h.handlers[msg.Topic]
		if !ok {
			log.Printf("No handler for topic %s", msg.Topic)
			session.MarkMessage(msg, "")
			continue
		}

		if err := handler(event); err != nil {
			log.Printf("Handler error for %s: %v", msg.Topic, err)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
