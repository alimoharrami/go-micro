package repository

import (
	"auth/internal/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) Create(ctx context.Context, permission *domain.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

func (r *PermissionRepository) GetByID(ctx context.Context, id uint) (*domain.Permission, error) {
	var permission domain.Permission
	if err := r.db.WithContext(ctx).First(&permission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) Update(ctx context.Context, permission *domain.Permission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

func (r *PermissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Permission{}, id).Error
}

func (r *PermissionRepository) List(ctx context.Context, offset, limit int) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&permissions).Error
	return permissions, err
}
