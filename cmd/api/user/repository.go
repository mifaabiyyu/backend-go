package user

import (
	"context"

	sqlc "github.com/mifaabiyyu/backend-go/internal/db/generated"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int64) (*sqlc.User, error)
	ListUsers(ctx context.Context, limit, offset int64) ([]sqlc.ListUsersRow, error)
}

type userRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) UserRepository {
	return &userRepository{q: q}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*sqlc.User, error) {
	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ListUsers(ctx context.Context, limit, offset int64) ([]sqlc.ListUsersRow, error) {
	return r.q.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
}
