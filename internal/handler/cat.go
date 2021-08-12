package handler

import (
	"net/http"

	"github.com/BorisMaslovskii/cats/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type Cat struct {
	Srv *service.CatService
}

func NewCat(srv *service.CatService) *Cat {
	return &Cat{Srv: srv}
}

func (h *Cat) GetById(c echo.Context) error {
	return nil
}

func (h *Cat) GetAll(c echo.Context) error {
	cats, err := h.Srv.GetAll(c.Request().Context())
	if err != nil {
		log.Errorf("GetAll %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, cats)

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
