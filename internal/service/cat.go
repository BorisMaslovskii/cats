package service

import (
	"context"

	"github.com/BorisMaslovskii/cats/internal/model"
	"github.com/BorisMaslovskii/cats/internal/repository"
)

type CatService struct {
	repo *repository.CatRepository
}

func NewCatService(rep *repository.CatRepository) *CatService {
	return &CatService{repo: rep}
}

func (c *CatService) GetAll(ctx context.Context) ([]*model.Cat, error) {
	return c.repo.GetAll(ctx)
}

func (c *CatService) GetByID(ctx context.Context, id string) (*model.Cat, error) {
	return c.repo.GetByID(ctx, id)
}

func (c *CatService) Create(ctx context.Context, cat *model.Cat) (int, error) {
	return c.repo.Create(ctx, cat)
}

func (c *CatService) Delete(ctx context.Context, id string) error {
	return c.repo.Delete(ctx, id)
}

func (c *CatService) Update(ctx context.Context, id string, cat *model.Cat) error {
	return c.repo.Update(ctx, id, cat)
}
