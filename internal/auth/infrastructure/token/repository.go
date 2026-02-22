package token

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{conn: pool}
}

func (repo *Repository) Create(ctx context.Context, token *RefreshToken) error {
	query, args, err := squirrel.Insert("refresh_tokens").
		Columns("id", "user_id", "access_id", "exp_at").
		Values(token.Id, token.UserId, token.AccessId, token.ExpAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "squirrel.ToSql")
	}

	_, err = repo.conn.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "repo.conn.Exec")
	}

	return nil
}

func (repo *Repository) GetByID(ctx context.Context, id uuid.UUID) (*RefreshToken, bool, error) {
	query, args, err := squirrel.Select("id", "user_id", "access_id", "exp_at").
		From("refresh_tokens").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, false, errors.Wrap(err, "squirrel.ToSql")
	}

	var token RefreshToken
	err = repo.conn.QueryRow(ctx, query, args...).Scan(&token.Id, &token.UserId, &token.AccessId, &token.ExpAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &token, false, nil
		}
		return nil, false, errors.Wrap(err, "repo.conn.QueryRow.Scan")
	}

	return &token, true, nil
}

func (repo *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := squirrel.Delete("refresh_tokens").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "squirrel.ToSql")
	}

	_, err = repo.conn.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "repo.conn.Exec")
	}

	return nil
}
