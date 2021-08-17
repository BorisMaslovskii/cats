package handler

import (
	"net/http"

	"github.com/BorisMaslovskii/cats/internal/model"
	"github.com/BorisMaslovskii/cats/internal/service"
	"github.com/golang-jwt/jwt"
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
	user := &model.User{
		Login:    userRec.Login,
		Password: userRec.Password,
	}

	tokenSignedString, err := h.AuthSrv.LogIn(c.Request().Context(), user)
	if err != nil {
		log.Errorf("User LogIn error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, tokenSignedString)
}

// JWTCheckAdmin middleware checks admin claims
func (h *AuthHandler) JWTCheckAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Get("user").(*jwt.Token)

		token, err := jwt.ParseWithClaims(tokenString.Raw, &service.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return nil, nil
		})

		claims, ok := token.Claims.(*service.JwtCustomClaims)
		if !ok {
			log.Error("JWTCheckAdmin JwtCustomClaims cast error")
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}

		isAdmin := claims.Admin
		if !isAdmin {
			log.Errorf("JWTCheckAdmin error at %v for user %v", c.Request().RequestURI, claims.Subject)
			return echo.NewHTTPError(http.StatusForbidden, "Access is forbidden. This api requires admin privileges")
		}
		return next(c)
	}
}
