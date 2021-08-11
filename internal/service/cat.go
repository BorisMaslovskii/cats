package service

import (
	"github.com/BorisMaslovskii/cats/internal/repository"
)

type CatService struct {
	rep *repository.Repo
}

func NewCatService(rep *repository.Repo) *CatService {
	return &CatService{rep: rep}
}
