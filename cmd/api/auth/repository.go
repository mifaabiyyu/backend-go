package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	sqlc "github.com/mifaabiyyu/backend-go/internal/db/generated"
)

type Repository interface {
	IsEmailExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.User, error)
}

type authRepo struct {
	q *sqlc.Queries
}

func NewAuthRepository(q *sqlc.Queries) Repository {
	return &authRepo{q}
}

func (r *authRepo) IsEmailExists(ctx context.Context, email string) (bool, error) {
	user, err := r.q.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return user.ID != 0, nil
}

func (r *authRepo) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	return r.q.CreateUser(ctx, arg)
}

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	user, err := r.q.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return sqlc.User{}, sql.ErrNoRows
		}
		return sqlc.User{}, err
	}
	return user, nil
}
