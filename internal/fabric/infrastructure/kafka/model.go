package kafka

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Items     []string
	Status    string
	CreatedAt time.Time
}

type OrderEvent struct {
	EventType string    `json:"event_type"`
	Payload   *Order    `json:"payload"`
	SentAt    time.Time `json:"sent_at"`
}
