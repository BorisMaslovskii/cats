package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/BorisMaslovskii/cats/internal/model"
	"github.com/BorisMaslovskii/cats/internal/service"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.
var hmacSampleSecret []byte

// UserRequest struct is used for binding the request content
type UserRequest struct {
	Login    string `query:"login" form:"login" json:"login"`
	Password string `query:"password" form:"password" json:"password"`
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

	// hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRec.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("User Create GenerateFromPassword error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := &model.User{
		Login:    userRec.Login,
		Password: string(hashedPassword),
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

// Delete handler func deletes a user
func (h *User) LogIn(c echo.Context) error {
	userRec := &UserRequest{}
	err := c.Bind(userRec)
	if err != nil {
		log.Errorf("User Update binding error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	reqUser := &model.User{
		Login:    userRec.Login,
		Password: userRec.Password,
	}

	loggedInUser, err := h.Srv.LogIn(c.Request().Context(), reqUser)
	if err != nil {
		log.Errorf("User LogIn error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(loggedInUser.Password), []byte(reqUser.Password))
	if err != nil {
		log.Errorf("User LogIn CompareHashAndPassword error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hmacSampleSecret = []byte("testSecret")

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.StandardClaims{
			Subject:   loggedInUser.Login,
			Issuer:    "cats project",
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		},
	)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		log.Errorf("User LogIn token.SignedString error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, tokenString)
}
