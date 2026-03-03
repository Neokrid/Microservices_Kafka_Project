package kafka

import (
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
