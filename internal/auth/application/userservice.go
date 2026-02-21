package auth

import (
	"context"
	"microservices_kafka_project/internal/auth/domain/dto/requests"
	"microservices_kafka_project/internal/auth/infrastructure/user"
)

type userService interface {
	CreateUser(ctx context.Context, credentials requests.RegisterCredentials) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string, password string) (*user.User, error)
}
