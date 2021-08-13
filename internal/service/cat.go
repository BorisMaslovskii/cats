// Package service provides usage methods for each service not depending on the type of repository
package service

import (
	"context"

	"github.com/BorisMaslovskii/cats/internal/model"
	"github.com/BorisMaslovskii/cats/internal/repository"
	"github.com/google/uuid"
)

// CatService struct
type CatService struct {
	repo repository.CatRepository
}

// NewCatService func creates new CatService
func NewCatService(rep repository.CatRepository) *CatService {
	return &CatService{repo: rep}
}

// GetAll func gets all cats
func (c *CatService) GetAll(ctx context.Context) ([]*model.Cat, error) {
	return c.repo.GetAll(ctx)
}

// GetByID func gets a cat by id
func (c *CatService) GetByID(ctx context.Context, id uuid.UUID) (*model.Cat, error) {
	return c.repo.GetByID(ctx, id)
}

// Create func creates a new cat
func (c *CatService) Create(ctx context.Context, cat *model.Cat) (uuid.UUID, error) {
	return c.repo.Create(ctx, cat)
}

// Delete func deletes a cat
func (c *CatService) Delete(ctx context.Context, id uuid.UUID) error {
	return c.repo.Delete(ctx, id)
}

// Update func updates a cat
func (c *CatService) Update(ctx context.Context, id uuid.UUID, cat *model.Cat) error {
	return c.repo.Update(ctx, id, cat)
}
