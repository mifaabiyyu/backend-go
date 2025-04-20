package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	// "github.com/mifaabiyyu/backend-go/api"
	"github.com/mifaabiyyu/backend-go/internal/auth"
	sqlc "github.com/mifaabiyyu/backend-go/internal/db/generated"
	"github.com/mifaabiyyu/backend-go/utils"
)

type userKey string

type AppAll struct {
	AppWrapper  *utils.AppWrapper
	Application *Application
}

const userCtx userKey = "user"

func (app *AppAll) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.AppWrapper.UnauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.AppWrapper.UnauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := app.Application.Authenticator.ValidateToken(token)
		if err != nil {
			app.AppWrapper.UnauthorizedErrorResponse(w, r, err)
			return
		}

		if !jwtToken.Valid {
			app.AppWrapper.UnauthorizedErrorResponse(w, r, fmt.Errorf("token invalid or expired"))
			return
		}

		claims, ok := jwtToken.Claims.(*auth.Claims)
		if !ok {
			app.AppWrapper.UnauthorizedErrorResponse(w, r, fmt.Errorf("invalid token claims"))
			return
		}

		userID := claims.UserID

		ctx := r.Context()

		user, err := app.getUser(ctx, userID)
		if err != nil {
			app.AppWrapper.UnauthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *AppAll) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.AppWrapper.UnauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.AppWrapper.UnauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.AppWrapper.UnauthorizedErrorResponse(w, r, err)
				return
			}

			username := app.Application.Config.Auth.Basic.User
			pass := app.Application.Config.Auth.Basic.Pass

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.AppWrapper.UnauthorizedErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *AppAll) getUser(ctx context.Context, userID int64) (*sqlc.User, error) {
	if !app.Application.Config.RedisCfg.Enabled {
		u, err := app.Application.Store.Queries.GetUserByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		return &u, nil
	}

	user, err := app.Application.CacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		u, err := app.Application.Store.Queries.GetUserByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		user = &u

		if err := app.Application.CacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (app *AppAll) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Application.Config.RateLimiter.Enabled {
			if allow, retryAfter := app.Application.RateLimiter.Allow(r.RemoteAddr); !allow {
				app.AppWrapper.RateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *AppAll) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userVal := r.Context().Value(userCtx)
			if userVal == nil {
				app.AppWrapper.UnauthorizedErrorResponse(w, r, fmt.Errorf("user not authenticated"))
				return
			}

			user, ok := userVal.(*sqlc.User)
			if !ok {
				app.AppWrapper.UnauthorizedErrorResponse(w, r, fmt.Errorf("invalid user type in context"))
				return
			}

			permissions, err := app.Application.Store.Queries.GetPermissionsByRoleID(r.Context(), user.RoleID.Int32)
			if err != nil {
				app.AppWrapper.ForbiddenResponse(w, r, fmt.Errorf("failed to retrieve permissions: %w", err))
				return
			}

			for _, p := range permissions {
				if p.Name == permission {
					next.ServeHTTP(w, r)
					return
				}
			}

			app.AppWrapper.ForbiddenResponse(w, r, fmt.Errorf("you don't have permission to perform this action"))
		})
	}
}

func (app *AppAll) RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(userCtx) == nil {
			app.AppWrapper.UnauthorizedErrorResponse(w, r, fmt.Errorf("unauthenticated"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
