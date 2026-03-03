package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic string, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

func (c *KafkaConsumer) Listen(ctx context.Context, handler func(context.Context, OrderEvent)) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("ошибка чтения из Kafka: %v", err)
			continue
		}

		log.Printf("DEBUG: Пришло сообщение из Kafka: %s", string(msg.Value))
		var event OrderEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Panicf("ошибка парсинга json: %v. Raw: %s", err, string(msg.Value))
			continue
		}

		handler(ctx, event)
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
