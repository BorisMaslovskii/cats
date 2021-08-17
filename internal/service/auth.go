package service

import (
	"context"
	"time"

	"github.com/BorisMaslovskii/cats/internal/config"
	"github.com/BorisMaslovskii/cats/internal/model"
	"github.com/BorisMaslovskii/cats/internal/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type JwtCustomClaims struct {
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

// UserService struct
type AuthService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

// NewUserService func creates new UserService
func NewAuthService(rep repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		repo: rep,
		cfg:  cfg,
	}
}

// LogIn func logins a user
func (s *AuthService) LogIn(ctx context.Context, reqUser *model.User) (tokenSignedString string, err error) {

	user, err := s.repo.GetByLogin(ctx, reqUser)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password))
	if err != nil {
		return "", err
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&JwtCustomClaims{
			user.Admin,
			jwt.StandardClaims{
				Subject:   user.Login,
				Issuer:    "cats-project",
				ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
			},
		},
	)

	hmacJWTSecret := []byte(s.cfg.HmacJWTSecretString)

	// Sign and get the complete encoded token as a string using the secret
	tokenSignedString, err = token.SignedString(hmacJWTSecret)
	if err != nil {
		return "", err
	}

	return tokenSignedString, nil
}
