package service

import (
	"auth/internal/domain"
	"auth/internal/repository"
	"context"
)

type RolePermissionService struct {
	repo *repository.RolePermissionRepository
}

func NewRolePermissionService(repo *repository.RolePermissionRepository) *RolePermissionService {
	return &RolePermissionService{repo: repo}
}

func (s *RolePermissionService) AssignPermissionToRole(ctx context.Context, roleID uint, permissionID uint) error {
	return s.repo.AssignPermissionToRole(ctx, roleID, permissionID)
}

func (s *RolePermissionService) RemovePermissionFromRole(ctx context.Context, roleID uint, permissionID uint) error {
	return s.repo.RemovePermissionFromRole(ctx, roleID, permissionID)
}

func (s *RolePermissionService) ListPermissionsByRole(ctx context.Context, roleID uint) ([]domain.Permission, error) {
	return s.repo.ListPermissionsByRole(ctx, roleID)
}
