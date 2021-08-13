package service

import (
	"context"

	"github.com/BorisMaslovskii/cats/internal/model"
	"github.com/BorisMaslovskii/cats/internal/repository"
	"github.com/google/uuid"
)

// CatService struct
type UserService struct {
	repo repository.UserRepository
}

// NewCatService func creates new CatService
func NewUserService(rep repository.UserRepository) *UserService {
	return &UserService{repo: rep}
}

// GetAll func gets all cats
func (c *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	return c.repo.GetAll(ctx)
}

// GetByID func gets a cat by id
func (c *UserService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return c.repo.GetByID(ctx, id)
}

// Create func creates a new cat
func (c *UserService) Create(ctx context.Context, cat *model.User) (uuid.UUID, error) {
	return c.repo.Create(ctx, cat)
}

// Delete func deletes a cat
func (c *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return c.repo.Delete(ctx, id)
}

// Update func updates a cat
func (c *UserService) Update(ctx context.Context, id uuid.UUID, cat *model.User) error {
	return c.repo.Update(ctx, id, cat)
}
