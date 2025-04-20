package auth

import (
	"github.com/mifaabiyyu/backend-go/internal/auth"
	"github.com/mifaabiyyu/backend-go/internal/store"
	"github.com/mifaabiyyu/backend-go/utils"
)

func InitAuthModule(store *store.Store, wrapper *utils.AppWrapper, authenticator auth.Authenticator) *Handler {
	repo := NewAuthRepository(store.Queries)

	service := NewAuthService(repo, authenticator)

	return &Handler{
		Service:    service,
		AppWrapper: wrapper,
	}
}
