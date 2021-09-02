package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/BorisMaslovskii/cats/internal/config"
	"github.com/BorisMaslovskii/cats/internal/handler"
	"github.com/BorisMaslovskii/cats/internal/repository"
	"github.com/BorisMaslovskii/cats/internal/service"
	"github.com/BorisMaslovskii/cats/internal/validator"
)

func main() {
	// we create catsSrv at the start to be able to choose between postgres and mongo repository for the CatService
	var catsSrv *service.CatService

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbPostgres, err := CreatePostgresDB(cfg)
	if err != nil {
		log.Errorf("CreatePostgresDB error %v", err)
		return
	}
	repoPostgres := repository.NewRepo(dbPostgres)

	// Choose DB for cats service
	if cfg.CatsDBType == "mongo" {
		mongoCollection, err := CreateMongoCollection(cfg)
		if err != nil {
			log.Errorf("CreateMongoCollection error %v", err)
			return
		}
		repoMongo := repository.NewRepoMongo(mongoCollection)

		catsSrv = service.NewCatService(repoMongo)

		log.Info("mongo DB is used")
	} else {
		catsSrv = service.NewCatService(repoPostgres)
		log.Info("postgres DB is used")
	}

	cats := handler.NewCatsHandler(catsSrv)

	userRepoPostgres := repository.NewUserRepo(dbPostgres)
	usersSrv := service.NewUserService(userRepoPostgres)
	users := handler.NewUsersHandler(usersSrv)

	authSrv := service.NewAuthService(userRepoPostgres, cfg)
	auth := handler.NewAuthHandler(authSrv)

	e := echo.New()
	e.Validator = validator.NewCustomValidator()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, this is a Cats service 3!")
	})

	// JWT
	hmacJWTSecret := []byte(cfg.HmacJWTSecretString)

	// Auth service
	e.GET("/auth/login", auth.LogIn)
	e.POST("/auth/login", auth.LogIn)

	// Cats service
	e.GET("/cats/:id", cats.GetByID, middleware.JWT(hmacJWTSecret))
	e.GET("/cats", cats.GetAll, middleware.JWT(hmacJWTSecret))
	e.POST("/cats", cats.Create, middleware.JWT(hmacJWTSecret))
	e.DELETE("/cats/:id", cats.Delete, middleware.JWT(hmacJWTSecret))
	e.PUT("/cats/:id", cats.Update, middleware.JWT(hmacJWTSecret))

	// Users service
	e.GET("/users/:id", users.GetByID, middleware.JWT(hmacJWTSecret), auth.JWTCheckAdmin)
	e.GET("/users", users.GetAll, middleware.JWT(hmacJWTSecret), auth.JWTCheckAdmin)
	e.POST("/users", users.Create, middleware.JWT(hmacJWTSecret), auth.JWTCheckAdmin)
	e.DELETE("/users/:id", users.Delete, middleware.JWT(hmacJWTSecret), auth.JWTCheckAdmin)
	e.PUT("/users/:id", users.Update, middleware.JWT(hmacJWTSecret), auth.JWTCheckAdmin)

	e.Logger.Fatal(e.Start(":" + cfg.HTTPPort))
}

// CreatePostgresDB func to simplify the main
func CreatePostgresDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host= %v port=%v user=%v password=%v sslmode=disable", cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateMongoCollection func to simplify the main
func CreateMongoCollection(cfg *config.Config) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%v:%v", cfg.MongoHost, cfg.MongoPort)))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	collection := client.Database(cfg.MongoDatabase).Collection(cfg.MongoCollection)

	return collection, nil
}
