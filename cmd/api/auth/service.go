package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mifaabiyyu/backend-go/internal/auth"
	sqlc "github.com/mifaabiyyu/backend-go/internal/db/generated"
	"github.com/mifaabiyyu/backend-go/internal/password"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*sqlc.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	repo          Repository
	authenticator auth.Authenticator
}

func NewAuthService(repo Repository, auth auth.Authenticator) Service {
	return &authService{
		repo:          repo,
		authenticator: auth,
	}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*sqlc.User, error) {
	exists, err := s.repo.IsEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, sqlc.CreateUserParams{
		Email:    req.Email,
		Password: hashedPassword,
		Username: req.Username,
		RoleID:   pgtype.Int4{Int32: 2, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password mismatch for email %s: %v", email, err)
		return "", errors.New("invalid email or password")
	}

	token, err := s.authenticator.GenerateToken(&auth.Claims{
		UserID: user.ID,
		RoleID: user.RoleID.Int32,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "backend-apps",
			Audience:  []string{"rahasia"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	})

	if err != nil {
		log.Printf("Error generating token for user %s: %v", email, err)
		return "", err
	}

	return token, nil
}
