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

// UserRequest struct is used for binding the request content
type UserRequest struct {
	Login    string `form:"login" json:"login"`
	Password string `form:"password" json:"password"`
}

// User handler struct provides handlers
type User struct {
	Srv *service.UserService
}

// NewUser func creates new User handler struct
func NewUser(srv *service.UserService) *User {
	return &User{Srv: srv}
}

// GetByID handler func gets a user by id
func (h *User) GetByID(c echo.Context) error {
	StringID := c.Param("id")
	id, err := uuid.Parse(StringID)
	if err != nil {
		log.Errorf("uuid FromBytes error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.Srv.GetByID(c.Request().Context(), id)
	if err != nil {
		log.Errorf("User GetById error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// GetAll handler func gets all users
func (h *User) GetAll(c echo.Context) error {
	users, err := h.Srv.GetAll(c.Request().Context())
	if err != nil {
		log.Errorf("User GetAll error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

// Create handler func creates a new user
func (h *User) Create(c echo.Context) error {
	userRec := &UserRequest{}
	err := c.Bind(userRec)
	if err != nil {
		log.Errorf("User Create binding error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user := &model.User{
		Login:    userRec.Login,
		Password: userRec.Password,
	}
	id, err := h.Srv.Create(c.Request().Context(), user)
	if err != nil {
		log.Errorf("User Create error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Created user № "+id.String())
}

// Update handler func updates a user
func (h *User) Update(c echo.Context) error {
	StringID := c.Param("id")
	id, err := uuid.Parse(StringID)
	if err != nil {
		log.Errorf("uuid FromBytes error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userRec := &UserRequest{}
	err = c.Bind(userRec)
	if err != nil {
		log.Errorf("User Update binding error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user := &model.User{
		Login:    userRec.Login,
		Password: userRec.Password,
	}
	err = h.Srv.Update(c.Request().Context(), id, user)
	if err != nil {
		log.Errorf("User Update error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Updated user № "+fmt.Sprint(id))
}

// Delete handler func deletes a user
func (h *User) Delete(c echo.Context) error {
	StringID := c.Param("id")
	id, err := uuid.Parse(StringID)
	if err != nil {
		log.Errorf("uuid FromBytes error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.Srv.Delete(c.Request().Context(), id)
	if err != nil {
		log.Errorf("User Delete error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, "Deleted user № "+id.String())
}
