package orders

import (
	"errors"
	"microservices_kafka_project/internal/orders/infrastructure/orders"
	"microservices_kafka_project/pkg/constants"
	"microservices_kafka_project/pkg/logger"
	"microservices_kafka_project/pkg/trx"
)

type Order struct {
	tx        trx.TransactionManager
	logger    logger.Logger
	orderRepo OrdersRepository
	order     orders.Order
}

func NewService(tx trx.TransactionManager, logger logger.Logger, ordersRepository OrdersRepository, order orders.Order) *Order {
	return &Order{
		tx:        tx,
		logger:    logger,
		orderRepo: ordersRepository,
		order:     order,
	}
}

func (o *Order) CanChangeStatus(newStatus string) error {
	if o.order.Status == constants.StatusDone {
		return errors.New("недопустимый  статус")
	}
	if o.order.Status == constants.StatusCreated {
		return errors.New("недопустимый  статус")
	}
	return nil
}
