package user

import sqlc "github.com/mifaabiyyu/backend-go/internal/db/generated"

func InitUserModule(queries *sqlc.Queries) *UserHandler {
	repo := NewUserRepository(queries)
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	return handler
}
