package domain

import (
	"errors"
	"time"
)

type Status string

const (
	StatusCreated    Status = "created"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Order struct {
	ID        string
	UserID    string
	Status    Status
	CreatedAt time.Time
}

func NewOrder(id, userID string) *Order {
	return &Order{
		ID:        id,
		UserID:    userID,
		Status:    StatusCreated,
		CreatedAt: time.Now(),
	}
}

func (o *Order) SetStatus(newStatus string) error {
	if newStatus == string(StatusDone) {
		return errors.New("connot change status, order is done")
	}
	o.Status = Status(newStatus)
	return nil
}
