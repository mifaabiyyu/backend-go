package user

import (
	"context"
)

type UserService interface {
	GetUser(ctx context.Context, id int64) (*User, error)
	ListUsers(ctx context.Context, limit, offset int64) ([]User, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUser(ctx context.Context, id int64) (*User, error) {
	u, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
		FullName: u.FullName,
	}, nil
}

func (s *userService) ListUsers(ctx context.Context, limit, offset int64) ([]User, error) {
	users, err := s.repo.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var result []User
	for _, u := range users {
		result = append(result, User{
			ID:       u.ID,
			Email:    u.Email,
			Username: u.Username,
			FullName: u.FullName,
		})
	}
	return result, nil
}
