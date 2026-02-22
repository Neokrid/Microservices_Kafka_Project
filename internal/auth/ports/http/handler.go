package authHttp

import (
	auth "microservices_kafka_project/internal/auth/application"
	"microservices_kafka_project/internal/auth/domain/dto/requests"
	"microservices_kafka_project/internal/auth/domain/services/token"

	"microservices_kafka_project/pkg/constants"
	"microservices_kafka_project/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *auth.AuthService
}

func NewAuthHandler(s *auth.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	ctx := c.Request.Context()
	var req requests.LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		_ = c.Error(errors.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	tokens, err := h.service.SignIn(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, LoginResponse{Token: tokens.Access})
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()
	var req requests.RegisterCredentials
	err := c.BindJSON(&req)
	if err != nil {
		_ = c.Error(errors.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	err = h.service.SignUp(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"Status": "ok",
	})
}

func (h *AuthHandler) RefreshTokens(c *gin.Context) {
	ctx := c.Request.Context()
	var req token.UserTokens
	err := c.BindJSON(&req)
	if err != nil {
		_ = c.Error(errors.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	tokens, err := h.service.RefreshTokens(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, LoginResponse{Token: tokens.Refresh})
}
