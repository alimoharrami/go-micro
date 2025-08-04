package service

import (
	"auth/internal/domain"
	"auth/internal/repository"
	"context"
)

type PermissionService struct {
	repo *repository.PermissionRepository
}

func NewPermissionService(repo *repository.PermissionRepository) *PermissionService {
	return &PermissionService{repo: repo}
}

func (s *PermissionService) Create(ctx context.Context, permission *domain.Permission) error {
	return s.repo.Create(ctx, permission)
}

func (s *PermissionService) GetByID(ctx context.Context, id uint) (*domain.Permission, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PermissionService) Update(ctx context.Context, permission *domain.Permission) error {
	return s.repo.Update(ctx, permission)
}

func (s *PermissionService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *PermissionService) List(ctx context.Context, offset, limit int) ([]domain.Permission, error) {
	return s.repo.List(ctx, offset, limit)
}
