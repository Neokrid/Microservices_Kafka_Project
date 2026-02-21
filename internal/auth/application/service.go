package auth

import (
	"context"
	"microservices_kafka_project/internal/auth/domain/dto/requests"
	"microservices_kafka_project/internal/auth/domain/services/token"
	auth "microservices_kafka_project/internal/auth/domain/services/user"
	"microservices_kafka_project/pkg/logger"
	"microservices_kafka_project/pkg/trx"
)

type AuthService struct {
	tx     trx.TransactionManager
	logger logger.Logger

	userService  userService
	tokenService tokenService
	repo         auth.UserRepo
}

func NewAuthService(
	r auth.UserRepo,
	tx trx.TransactionManager,
	logger logger.Logger,
	userService userService,
	tokenService tokenService,
) *AuthService {
	return &AuthService{
		repo:         r,
		tx:           tx,
		logger:       logger,
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *AuthService) SignUp(ctx context.Context, credentials requests.RegisterCredentials) error {
	user, err := s.userService.CreateUser(ctx, credentials)
	if err != nil {
		return err
	}
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) SignIn(ctx context.Context, req requests.LoginRequest) (*token.UserTokens, error) {
	u, err := s.userService.GetUserByEmail(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return s.tokenService.GenerateUserTokens(ctx, u.Id)

}

func (s *AuthService) RefreshTokens(ctx context.Context, req token.UserTokens) (*token.UserTokens, error) {
	return s.tokenService.RefreshTokens(ctx, req.Access, req.Refresh)
}
