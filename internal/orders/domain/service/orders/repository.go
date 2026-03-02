package orders

import (
	"context"
	"microservices_kafka_project/internal/orders/infrastructure/orders"

	"github.com/google/uuid"
)

type OrdersRepository interface {
	Save(ctx context.Context, order orders.Order) error
	GetByID(ctx context.Context, orderId uuid.UUID) (*orders.Order, error)
	GetByUserID(ctx context.Context, userId uuid.UUID) ([]*orders.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type OrdersPublisher interface {
	PublishOrderCreated(ctx context.Context, order *orders.Order) error
	PublishStatusUpdated(ctx context.Context, order *orders.Order) error
}
