package user

import (
	"context"
	"database/sql"
	"microservices_kafka_project/internal/common"
	customErrors "microservices_kafka_project/internal/common/customErrors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{conn: pool}
}

func (repo *Repository) CreateUser(ctx context.Context, item *User) error {
	query := `
		INSERT INTO users VALUES
		($1, $2, $3, $4, $5)
	`
	_, err := repo.conn.Exec(ctx, query, item.Id, item.Username, item.Email, item.Password, item.CreatedAt)
	if err != nil {
		if common.IsUniqueErr(err) {
			return customErrors.NotUnique
		}
	}
	return err
}

func (repo *Repository) GetUser(ctx context.Context, filter UserFilter) (*User, bool, error) {
	var output User
	builder := squirrel.Select("u.*").From("users u")

	if filter.Id != nil {
		builder = builder.Where(squirrel.Eq{"id": filter.Id})
	}

	if filter.Email != nil {
		builder = builder.Where(squirrel.Eq{"email": filter.Email})
	}

	if filter.Limit > 0 {
		builder.Limit(filter.Limit)
	}
	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, false, errors.Wrap(err, "squirrel.ToSql")
	}
	err = repo.conn.QueryRow(ctx, query, args...).
		Scan(&output.Id, &output.Username, &output.Email, &output.Password, &output.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &output, false, nil
		}
		if common.IsUniqueErr(err) {
			return &output, false, customErrors.NotUnique
		}
		return &output, false, errors.Wrap(err, "repo.conn.QueryRow.Scan")
	}

	return &output, true, nil
}
