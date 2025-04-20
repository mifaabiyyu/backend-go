package store

import (
	"github.com/jackc/pgx/v5/pgxpool"
	sqlc "github.com/mifaabiyyu/backend-go/internal/db/generated"
)

type Store struct {
	Queries *sqlc.Queries
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Queries: sqlc.New(pool),
	}
}
