package repository

import (
	"context"
	"errors"
	"go-blog/internal/domain"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *PostRepository) GetByID(ctx context.Context, id uint) (*domain.Post, error) {
	var post domain.Post
	if err := r.db.WithContext(ctx).First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *PostRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Post{}, id).Error
}

func (r *PostRepository) List(ctx context.Context, offset, limit int) ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.WithContext(ctx).Find(&posts).Error
	// .Offset(offset).Limit(limit)

	return posts, err
}
