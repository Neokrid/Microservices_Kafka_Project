package orders

import (
	"context"
	"errors"
	"log"
	"microservices_kafka_project/internal/orders/domain/service/orders"
	orderDTO "microservices_kafka_project/internal/orders/infrastructure/orders"
	"microservices_kafka_project/pkg/constants"

	"github.com/google/uuid"
)

type OrderService struct {
	repo      orders.OrdersRepository
	publisher orders.OrdersPublisher
	order     OrdersService
}

func NewOrdersService(r orders.OrdersRepository, p orders.OrdersPublisher, o OrdersService) *OrderService {
	return &OrderService{
		repo:      r,
		publisher: p,
		order:     o,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uuid.UUID, items []string) (*orderDTO.Order, error) {
	order := orderDTO.Order{
		ID:     uuid.New(),
		UserID: userID,
		Items:  items,
		Status: constants.StatusCreated,
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return nil, err
	}
	if err := s.publisher.PublishOrderCreated(ctx, &order); err != nil {

		log.Printf("Критическая ошибка Kafka: %v", err)
		return nil, nil
	}
	return &order, nil
}

func (s *OrderService) GetUserOrder(ctx context.Context, userID, orderID uuid.UUID) (*orderDTO.Order, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("это не ваш заказ!")
	}

	return order, nil
}

func (s *OrderService) GetAllUserOrders(ctx context.Context, userID uuid.UUID) ([]*orderDTO.Order, error) {
	orders, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error {

	if err := s.order.CanChangeStatus(status); err != nil {
		return err
	}

	if err := s.repo.UpdateStatus(ctx, orderID, status); err != nil {
		return err
	}

	updatedOrder, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		log.Printf("Ошибка при получении заказа для Kafka: %v", err)
		return err
	}
	if err := s.publisher.PublishStatusUpdated(ctx, updatedOrder); err != nil {

		log.Printf("Критическая ошибка Kafka: %v", err)
		return err
	}
	return nil
}
