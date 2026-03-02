package publichttp

import (
	orders "microservices_kafka_project/internal/orders/application"
	"microservices_kafka_project/pkg/constants"
	"microservices_kafka_project/pkg/errors"
	"microservices_kafka_project/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PublicOrderHandler struct {
	service *orders.OrderService
}

func NewPublicOrderHandler(s *orders.OrderService) *PublicOrderHandler {
	return &PublicOrderHandler{
		service: s,
	}
}

func (h *PublicOrderHandler) CreateOrder(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := utils.GetUserId(ctx)
	if err != nil {
		_ = c.Error(errors.NewBadRequestError(err.Error(), constants.UserIdTypeMismatch))
		return
	}
	var req struct {
		Items []string `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), userID, req.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *PublicOrderHandler) GetOrderById(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := utils.GetUserId(ctx)
	if err != nil {
		_ = c.Error(errors.NewBadRequestError(err.Error(), constants.UserIdTypeMismatch))
		return
	}
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	order, err := h.service.GetUserOrder(c.Request.Context(), userID, orderID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *PublicOrderHandler) GetAllUserOrders(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := utils.GetUserId(ctx)
	if err != nil {
		_ = c.Error(errors.NewBadRequestError(err.Error(), constants.UserIdTypeMismatch))
		return
	}
	orders, err := h.service.GetAllUserOrders(ctx, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}
