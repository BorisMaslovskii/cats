package handler

import (
	"net/http"

	"github.com/BorisMaslovskii/cats/internal/model"
	"github.com/BorisMaslovskii/cats/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// UsersHandler handler struct provides handlers
type AuthHandler struct {
	AuthSrv *service.AuthService
}

// NewUsersHandler func creates new User handler struct
func NewAuthHandler(authSrv *service.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthSrv: authSrv,
	}
}

// LogIn handler func loggs in a user
func (h *AuthHandler) LogIn(c echo.Context) error {
	userRec := &UserRequest{}
	err := c.Bind(userRec)
	if err != nil {
		log.Errorf("User LogIn binding error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	reqUser := &model.User{
		Login:    userRec.Login,
		Password: userRec.Password,
	}

	tokenSignedString, err := h.AuthSrv.LogIn(c.Request().Context(), reqUser)
	if err != nil {
		log.Errorf("User LogIn error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, tokenSignedString)
}
