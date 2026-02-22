package auth

import (
	"context"
	"microservices_kafka_project/internal/auth/infrastructure/user"
)

type UserRepo interface {
	CreateUser(ctx context.Context, item *user.User) error
	GetUser(ctx context.Context, filter user.UserFilter) (*user.User, bool, error)
}
