package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"microservices_kafka_project/internal/orders/infrastructure/orders"
	"time"

	"github.com/segmentio/kafka-go"
)

type OrderEvent struct {
	EventType string        `json:"event_type"`
	Payload   *orders.Order `json:"payload"`
	SentAt    time.Time     `json:"sent_at"`
}

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {
	fmt.Printf("KAFKA DEBUG: Подключаемся к %v, топик: '%s'\n", brokers, topic)
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			Async:                  false,
			Logger:                 kafka.LoggerFunc(log.Printf),
			ErrorLogger:            kafka.LoggerFunc(log.Printf),
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *KafkaPublisher) PublishOrderCreated(ctx context.Context, order *orders.Order) error {
	event := OrderEvent{
		EventType: "OrderCreated",
		Payload:   order,
		SentAt:    time.Now(),
	}

	return p.send(ctx, order.ID.String(), event)
}

func (p *KafkaPublisher) PublishStatusUpdated(ctx context.Context, order *orders.Order) error {
	event := OrderEvent{
		EventType: "StatusUpdated",
		Payload:   order,
		SentAt:    time.Now(),
	}

	return p.send(ctx, order.ID.String(), event)
}

func (p *KafkaPublisher) send(ctx context.Context, key string, event OrderEvent) error {
	msgBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: msgBytes,
	})

	if err != nil {
		return fmt.Errorf("kafka write: %w", err)
	}

	return nil
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
