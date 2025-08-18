package service

import (
	"context"
	"errors"
	"fmt"
	"go-blog/internal/domain"
	"go-blog/internal/repository"
	"math"
)

type PostService struct {
	repo *repository.PostRepository
}

// CreateUserInput defines user creation fields.
type CreatePostInput struct {
	Title   string
	Content string
}

// UpdateUserInput defines fields allowed for update.
type UpdatePostInput struct {
	Title   *string
	Content *string
}

// NewUserService initializes UserService.
func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

// GetByID fetches a user by ID
func (s *PostService) GetByID(ctx context.Context, id uint) (*domain.Post, error) {
	return s.repo.GetByID(ctx, id)
}

// Create hashes password and creates a new user.
func (s *PostService) Create(ctx context.Context, input CreatePostInput) (*domain.Post, error) {

	post := &domain.Post{
		Title:   input.Title,
		Content: input.Content,
	}

	if err := s.repo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

func (s *PostService) Update(ctx context.Context, id uint, input UpdatePostInput) (*domain.Post, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if post == nil {
		return nil, errors.New("user not found")
	}

	// Apply updates only if fields are provided
	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Content != nil {
		post.Content = *input.Content
	}

	if err := s.repo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

func (s *PostService) List(ctx context.Context, page, limit int) (map[string]interface{}, error) {
	offset := (page - 1) * limit
	posts, total, err := s.repo.List(ctx, offset, limit)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}

	result := map[string]interface{}{
		"data": posts,
		"pagination": map[string]interface{}{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	return result, nil
}

func (s *PostService) Delete(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("there is problem in deleting post. %w", err)
	}

	return nil
}

//todo get posts paginate
