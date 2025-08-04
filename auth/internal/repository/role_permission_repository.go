package repository

import (
	"auth/internal/domain"
	"context"

	"gorm.io/gorm"
)

type RolePermissionRepository struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) *RolePermissionRepository {
	return &RolePermissionRepository{db: db}
}

func (r *RolePermissionRepository) AssignPermissionToRole(ctx context.Context, roleID uint, permissionID uint) error {
	rp := domain.RolePermission{RoleID: roleID, PermissionID: permissionID}
	return r.db.WithContext(ctx).Create(&rp).Error
}

func (r *RolePermissionRepository) RemovePermissionFromRole(ctx context.Context, roleID uint, permissionID uint) error {
	return r.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&domain.RolePermission{}).Error
}

func (r *RolePermissionRepository) ListPermissionsByRole(ctx context.Context, roleID uint) ([]domain.Permission, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).Preload("Permissions").First(&role, roleID).Error
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}
