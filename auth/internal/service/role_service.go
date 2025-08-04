package service

import (
	"auth/internal/domain"
	"auth/internal/repository"
	"context"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) Create(ctx context.Context, role *domain.Role) error {
	return s.repo.Create(ctx, role)
}

func (s *RoleService) GetByID(ctx context.Context, id uint) (*domain.Role, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RoleService) Update(ctx context.Context, role *domain.Role) error {
	return s.repo.Update(ctx, role)
}

func (s *RoleService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *RoleService) List(ctx context.Context, offset, limit int) ([]domain.Role, error) {
	return s.repo.List(ctx, offset, limit)
}
