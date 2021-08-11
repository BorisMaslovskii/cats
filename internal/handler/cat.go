package protocol

import (
	"github.com/BorisMaslovskii/cats/internal/service"
	"github.com/labstack/echo/v4"
)

type Cat struct {
	srv *service.CatService
}

func NewCat(srv *service.CatService) *Cat {
	return &Cat{srv: srv}
}

func (h *Cat) Get(c echo.Context) error {
	return nil
}

func (h *Cat) GetAll(c echo.Context) error {
	return nil
}

func (h *Cat) Create(c echo.Context) error {
	return nil
}

func (h *Cat) Update(c echo.Context) error {
	return nil
}

func (h *Cat) Delete(c echo.Context) error {
	return nil
}
