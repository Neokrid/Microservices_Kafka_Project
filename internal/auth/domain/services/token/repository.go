package token

import (
	"context"
	"microservices_kafka_project/internal/auth/infrastructure/token"
	"time"

	"github.com/google/uuid"
)

type tokenRepo interface {
	Create(ctx context.Context, token *token.RefreshToken) error
	GetByID(ctx context.Context, id uuid.UUID) (*token.RefreshToken, bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context, cutoffTime time.Time) error
}
