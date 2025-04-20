package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	sqlc "github.com/mifaabiyyu/backend-go/internal/db/generated"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*sqlc.User, error)
		Set(context.Context, *sqlc.User) error
		Delete(context.Context, int64)
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rbd},
	}
}
