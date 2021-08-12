package handler

import (
	"fmt"
	"net/http"

	"github.com/BorisMaslovskii/cats/internal/model"
	"github.com/BorisMaslovskii/cats/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type CatRequest struct {
	Name  string `form:"name" json:"name"`
	Color string `form:"color" json:"color"`
}

type Cat struct {
	Srv *service.CatService
}

func NewCat(srv *service.CatService) *Cat {
	return &Cat{Srv: srv}
}

func (h *Cat) GetById(c echo.Context) error {
	id := c.Param("id")
	cat, err := h.Srv.GetByID(c.Request().Context(), id)
	if err != nil {
		log.Errorf("Cat GetById %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, cat)
}

func (h *Cat) GetAll(c echo.Context) error {
	cats, err := h.Srv.GetAll(c.Request().Context())
	if err != nil {
		log.Errorf("Cat GetAll %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, cats)
}

func (h *Cat) Create(c echo.Context) error {
	catRec := &CatRequest{}
	err := c.Bind(catRec)
	if err != nil {
		log.Errorf("Cat Create binding %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	cat := &model.Cat{
		Name:  catRec.Name,
		Color: catRec.Color,
	}
	catID, err := h.Srv.Create(c.Request().Context(), cat)
	if err != nil {
		log.Errorf("Cat Create %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.String(http.StatusOK, "Created cat № "+fmt.Sprint(catID))
}

func (h *Cat) Update(c echo.Context) error {
	id := c.Param("id")
	catRec := &CatRequest{}
	err := c.Bind(catRec)
	if err != nil {
		log.Errorf("Cat Update binding %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	cat := &model.Cat{
		Name:  catRec.Name,
		Color: catRec.Color,
	}
	err = h.Srv.Update(c.Request().Context(), id, cat)
	if err != nil {
		log.Errorf("Cat Update %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.String(http.StatusOK, "Updated cat № "+fmt.Sprint(id))
}

func (h *Cat) Delete(c echo.Context) error {
	id := c.Param("id")
	err := h.Srv.Delete(c.Request().Context(), id)
	if err != nil {
		log.Errorf("Cat Delete %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.String(http.StatusOK, "Deleted cat № "+fmt.Sprint(id))
}
