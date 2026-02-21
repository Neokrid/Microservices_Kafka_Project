package user

import (
	"context"
	"microservices_kafka_project/internal/common"
	errors "microservices_kafka_project/internal/common/customErrors"
	"microservices_kafka_project/pkg/database/postgres"
)

type Repository struct {
	conn postgres.Connection
}

func NewRepository(conn postgres.Connection) *Repository {
	return &Repository{conn: conn}
}

func (repo *Repository) CreateUser(ctx context.Context, item *User) error {
	query := `
		INSERT INTO users VALUES
		($1, $2, $3, $4, $5)
	`
	_, err := repo.conn.Exec(ctx, query, item.Id, item.Username, item.Email, item.Password, item.CreatedAt)
	if err != nil {
		if common.IsUniqueErr(err) {
			return errors.NotUnique
		}
	}
	return err
}
