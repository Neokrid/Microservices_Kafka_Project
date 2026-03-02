package token

import (
	"context"
	"microservices_kafka_project/internal/auth/infrastructure/token"
	"microservices_kafka_project/pkg/utils"
	"time"

	customError "microservices_kafka_project/internal/common/customErrors"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Service struct {
	refreshTokenTTL time.Duration
	accessTokenTTL  time.Duration
	secret          string
	tokenRepo       tokenRepo
}

func NewService(
	refreshTokenTTL time.Duration,
	accessTokenTTL time.Duration,
	secret string,
	tokenRepo tokenRepo,
) *Service {
	a := refreshTokenTTL.Seconds()
	b := accessTokenTTL.Seconds()
	_, _ = a, b

	return &Service{
		refreshTokenTTL: refreshTokenTTL,
		accessTokenTTL:  accessTokenTTL,
		secret:          secret,
		tokenRepo:       tokenRepo,
	}
}

func (s *Service) CreateUserTokens(id uuid.UUID) (*UserTokens, uuid.UUID, uuid.UUID, error) {
	jtiAccess := uuid.New()
	jtiRefresh := uuid.New()
	access, err := generateToken(jtiAccess, id, s.accessTokenTTL, s.secret)
	if err != nil {
		return nil, uuid.UUID{}, uuid.UUID{}, err
	}
	refresh, err := generateToken(jtiRefresh, id, s.refreshTokenTTL, s.secret)
	return &UserTokens{Access: access, Refresh: refresh}, jtiAccess, jtiRefresh, nil
}

func (s *Service) ParseToken(token string) (*CustomClaims, error) {
	return parseToken(s.secret, token)
}

func (s *Service) GenerateUserTokens(ctx context.Context, userId uuid.UUID) (*UserTokens, error) {
	t, accessId, refreshId, err := s.CreateUserTokens(userId)
	if err != nil {
		return nil, errors.Wrap(err, ".GenerateUserTokens")
	}

	return t, s.tokenRepo.Create(ctx, &token.RefreshToken{
		Id:       refreshId,
		UserId:   userId,
		AccessId: accessId,
		ExpAt:    utils.GetCurrentUTCTime().Add(s.refreshTokenTTL),
	})

}

func (s *Service) RefreshTokens(ctx context.Context, access, refresh string) (*UserTokens, error) {
	aToken, err := parseToken(s.secret, access)
	if err != nil {
		return nil, err
	}
	rToken, err := parseToken(s.secret, refresh)
	if err != nil {
		return nil, err
	}
	tokenId, err := uuid.Parse(rToken.ID)
	if err != nil {
		return nil, customError.TokenClaimsError
	}
	dbToken, ex, err := s.tokenRepo.GetByID(ctx, tokenId)
	if err != nil {
		return nil, err
	}
	if !ex {
		return nil, customError.TokenDontExist
	}
	if dbToken.AccessId.String() != aToken.ID ||
		dbToken.Id != tokenId ||
		dbToken.UserId != aToken.UserId {
		return nil, customError.TokensDontMatch
	}
	t, err := s.GenerateUserTokens(ctx, aToken.UserId)
	if err != nil {
		return nil, err
	}
	
	err = s.tokenRepo.Delete(ctx, dbToken.Id)
	if err != nil {
		return nil, err
	}
	return t, nil
}
