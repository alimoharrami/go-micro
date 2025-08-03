package service

import (
	"context"
	"errors"
	"fmt"
	"go-blog/internal/domain"
	"go-blog/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// CreateUserInput defines user creation fields.
type LoginInput struct {
	Email    string
	Password string
}

func (l *AuthService) Login(ctx context.Context, input LoginInput) (*domain.User, *string, error) {
	user, err := l.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to login: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, nil, errors.New("invalid email or password")
	}

	token, err := Generate(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, &token, nil
}
