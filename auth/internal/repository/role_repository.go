package repository

import (
	"auth/internal/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *RoleRepository) GetByID(ctx context.Context, id uint) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *RoleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Role{}, id).Error
}

func (r *RoleRepository) List(ctx context.Context, offset, limit int) ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&roles).Error
	return roles, err
}
