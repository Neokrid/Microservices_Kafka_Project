package auth

import (
	"context"
	"microservices_kafka_project/internal/auth/domain/services/token"

	"github.com/google/uuid"
)

type tokenService interface {
	GenerateUserTokens(ctx context.Context, id uuid.UUID) (*token.UserTokens, error)
	ParseToken(token string) (*token.CustomClaims, error)
	RefreshTokens(ctx context.Context, access, refresh string) (*token.UserTokens, error)
}
