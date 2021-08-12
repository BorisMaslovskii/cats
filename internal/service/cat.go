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
