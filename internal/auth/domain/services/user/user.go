package auth

import (
	"context"
	"crypto/sha512"
	"fmt"
	"microservices_kafka_project/internal/auth/domain/dto/requests"
	"microservices_kafka_project/internal/auth/infrastructure/user"
	errors "microservices_kafka_project/internal/common/customErrors"
	"microservices_kafka_project/pkg/logger"
	"microservices_kafka_project/pkg/trx"

	"github.com/google/uuid"
)

var salt string = "NiceDick"

type Service struct {
	tx       trx.TransactionManager
	logger   logger.Logger
	userRepo UserRepo
}

func NewService(tx trx.TransactionManager, logger logger.Logger, userRepo UserRepo) *Service {
	return &Service{
		tx:       tx,
		logger:   logger,
		userRepo: userRepo,
	}
}

func (s *Service) CreateUser(ctx context.Context, credentials requests.RegisterCredentials) (*user.User, error) {
	user := user.User{
		Id:       uuid.New(),
		Username: credentials.Username,
		Email:    credentials.Email,
		Password: generatePasswordHash(credentials.Password),
	}
	return &user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string, password string) (*user.User, error) {
	targetEntityUser, ex, err := s.userRepo.GetUser(ctx, user.UserFilter{
		Email: &email,
	})
	if err != nil {
		return nil, err
	}

	if !ex {
		return nil, errors.UserNotFound
	}

	if password != "" {
		if !s.comparePassword(password, targetEntityUser.Password) {
			return nil, errors.IncorrectPassword
		}
	}
	return targetEntityUser, nil
}

func generatePasswordHash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *Service) comparePassword(origin, existed string) bool {
	return generatePasswordHash(origin) == existed
}
