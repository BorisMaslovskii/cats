//nolint:dupl // For now, we want a more readable code than optimized
// Package handler provides handler structs and methods for each service
package handler

import (
	"fmt"
	"net/http"

	"github.com/BorisMaslovskii/cats/internal/model"
	"github.com/BorisMaslovskii/cats/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// CatsHandler handler struct provides handlers
type CatsHandler struct {
	Srv *service.CatService
}

// NewCatsHandler func creates new Cat handler struct
func NewCatsHandler(srv *service.CatService) *CatsHandler {
	return &CatsHandler{Srv: srv}
}

// GetByID handler func gets a cat by id
func (h *CatsHandler) GetByID(c echo.Context) error {
	StringID := c.Param("id")
	id, err := uuid.Parse(StringID)
	if err != nil {
		log.Errorf("uuid FromBytes error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	cat, err := h.Srv.GetByID(c.Request().Context(), id)
	if err != nil {
		log.Errorf("Cat GetById error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, cat)
}

// GetAll handler func gets all cats
func (h *CatsHandler) GetAll(c echo.Context) error {
	cats, err := h.Srv.GetAll(c.Request().Context())
	if err != nil {
		log.Errorf("Cat GetAll error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, cats)
}

// Create handler func creates a new cat
func (h *CatsHandler) Create(c echo.Context) error {
	catRec := &model.CatRequest{}
	err := c.Bind(catRec)
	if err != nil {
		log.Errorf("Cat Create binding error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	cat := &model.Cat{
		Name:  catRec.Name,
		Color: catRec.Color,
	}
	id, err := h.Srv.Create(c.Request().Context(), cat)
	if err != nil {
		log.Errorf("Cat Create error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Created cat № "+id.String())
}

// Update handler func updates a cat
func (h *CatsHandler) Update(c echo.Context) error {
	StringID := c.Param("id")
	id, err := uuid.Parse(StringID)
	if err != nil {
		log.Errorf("uuid FromBytes error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	catRec := &model.CatRequest{}
	err = c.Bind(catRec)
	if err != nil {
		log.Errorf("Cat Update binding error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	cat := &model.Cat{
		Name:  catRec.Name,
		Color: catRec.Color,
	}
	err = h.Srv.Update(c.Request().Context(), id, cat)
	if err != nil {
		log.Errorf("Cat Update error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Updated cat № "+fmt.Sprint(id))
}

// Delete handler func deletes a cat
func (h *CatsHandler) Delete(c echo.Context) error {
	StringID := c.Param("id")
	id, err := uuid.Parse(StringID)
	if err != nil {
		log.Errorf("uuid FromBytes error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.Srv.Delete(c.Request().Context(), id)
	if err != nil {
		log.Errorf("Cat Delete error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Deleted cat № "+id.String())
}
