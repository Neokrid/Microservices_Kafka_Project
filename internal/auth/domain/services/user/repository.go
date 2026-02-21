package auth

import (
	"context"
	"microservices_kafka_project/internal/auth/infrastructure/user"

	"github.com/google/uuid"
)

type UserRepo interface {
	CreateUser(ctx context.Context, item *user.User) error
	GetUser(ctx context.Context, filter user.UserFilter) (*user.User, bool, error)
	UpdateUser(ctx context.Context, userId uuid.UUID, updateParams *user.UserUpdateParams) error
}
