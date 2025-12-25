package service

import (
	"context"
	"errors"
	"fmt"
	"go-blog/internal/domain"
	"go-blog/internal/repository"
	"math"

	"go-blog/external/protos/userpb"
)

type PostService struct {
	repo       *repository.PostRepository
	userClient userpb.UserServiceClient
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
func NewPostService(repo *repository.PostRepository, userClient userpb.UserServiceClient) *PostService {
	return &PostService{repo: repo, userClient: userClient}
}

// GetByID fetches a user by ID
func (s *PostService) GetByID(ctx context.Context, id uint) (map[string]interface{}, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	// Fetch user
	resp, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: fmt.Sprintf("%d", post.AuthorID)})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	result := map[string]interface{}{
		"post": post,
		"user": resp.User,
	}

	return result, nil
}

// Create hashes password and creates a new user.
func (s *PostService) Create(ctx context.Context, userId uint, input CreatePostInput) (*domain.Post, error) {
	// Validate user exists
	resp, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: fmt.Sprintf("%d", userId)})
	if err != nil {
		return nil, fmt.Errorf("failed to validate user: %w", err)
	}
	if resp.User == nil {
		return nil, fmt.Errorf("user not found")
	}

	post := &domain.Post{
		Title:   input.Title,
		Content: input.Content,
		AuthorID: userId,
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

	// Collect unique author IDs
	authorIDs := make(map[string]bool)
	ids := []string{}
	for _, post := range posts {
		idStr := fmt.Sprintf("%d", post.AuthorID)
		if !authorIDs[idStr] {
			authorIDs[idStr] = true
			ids = append(ids, idStr)
		}
	}

	// Fetch users
	var users map[string]*userpb.User
	if len(ids) > 0 {
		resp, err := s.userClient.GetUsersByIDs(ctx, &userpb.GetUsersRequest{Ids: ids})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch users: %w", err)
		}
		users = make(map[string]*userpb.User)
		for _, user := range resp.Users {
			users[user.Id] = user
		}
	}

	// Attach users to posts
	postsWithUsers := []map[string]interface{}{}
	for _, post := range posts {
		postMap := map[string]interface{}{
			"id":        post.ID,
			"title":     post.Title,
			"content":   post.Content,
			"author_id": post.AuthorID,
		}
		if user, ok := users[fmt.Sprintf("%d", post.AuthorID)]; ok {
			postMap["author"] = user
		}
		postsWithUsers = append(postsWithUsers, postMap)
	}

	result := map[string]interface{}{
		"data": postsWithUsers,
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
