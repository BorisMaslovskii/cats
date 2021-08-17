// Package service comment
//nolint:dupl // For now, we want a more readable code than optimized
package service

import (
	"context"

	"github.com/BorisMaslovskii/cats/internal/model"
	"github.com/BorisMaslovskii/cats/internal/repository"
	"github.com/google/uuid"
)

// UserService struct
type UserService struct {
	repo repository.UserRepository
}

// NewUserService func creates new UserService
func NewUserService(rep repository.UserRepository) *UserService {
	return &UserService{repo: rep}
}

// GetAll func gets all users
func (s *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	return s.repo.GetAll(ctx)
}

// GetByID func gets a user by id
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

// Create func creates a new user
func (s *UserService) Create(ctx context.Context, cat *model.User) (uuid.UUID, error) {
	return s.repo.Create(ctx, cat)
}

// Delete func deletes a user
func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// Update func updates a user
func (s *UserService) Update(ctx context.Context, id uuid.UUID, cat *model.User) error {
	return s.repo.Update(ctx, id, cat)
}
