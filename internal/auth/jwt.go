package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

type Claims struct {
	UserID int64 `json:"user_id"`
	RoleID int32 `json:"role_id"`
	jwt.RegisteredClaims
}

// Implementasi interface `jwt.Claims`
func (c Claims) Valid() error {
	return nil
}

func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{secret, iss, aud}
}

func (j *JWTAuthenticator) GenerateToken(c jwt.Claims) (string, error) {
	claims, ok := c.(*Claims)
	if !ok {
		return "", fmt.Errorf("invalid claims type")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
}
