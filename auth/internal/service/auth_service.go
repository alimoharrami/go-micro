package service

import (
	"auth/internal/domain"
	"auth/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/alimoharrami/go-micro/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo                  *repository.UserRepository
	rolePermissionService *RolePermissionService
	roleService           *RoleService
	permissionService     *PermissionService
}

func NewAuthService(
	repo *repository.UserRepository,
	rolePermissionService *RolePermissionService,
	roleService *RoleService,
	permissionService *PermissionService,
) *AuthService {
	return &AuthService{
		repo:                  repo,
		rolePermissionService: rolePermissionService,
		roleService:           roleService,
		permissionService:     permissionService,
	}
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

	role, err := l.roleService.GetByID(ctx, user.RoleID)
	if err != nil {
		return nil, nil, err
	}

	permissions, err := l.rolePermissionService.ListPermissionsByRole(ctx, role.ID)
	if err != nil {
		return nil, nil, err
	}

	var permissionNames []string
	for _, p := range permissions {
		permissionNames = append(permissionNames, p.Key)
	}

	token, err := auth.Generate(user.ID, []string{role.Name}, permissionNames)
	if err != nil {
		return nil, nil, err
	}

	return user, &token, nil
}
