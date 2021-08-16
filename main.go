package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BorisMaslovskii/cats/internal/config"
	"github.com/BorisMaslovskii/cats/internal/handler"
	"github.com/BorisMaslovskii/cats/internal/repository"
	"github.com/BorisMaslovskii/cats/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// we create srv at the start to be able to choose between postgres and mongo repository for the CatService
	var srv *service.CatService

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	connStr := fmt.Sprintf("user=%v password=%v sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Errorf("sql Open error %v", err)
		return
	}
	defer func() {
		errd := db.Close()
		if errd != nil {
			log.Errorf("defer db Close error %v", errd)
		}
	}()

	err = db.Ping()
	if err != nil {
		log.Errorf("db Ping error %v", err)
		return
	}

	repoPostgres := repository.NewRepo(db)

	DBType := ""
	if len(os.Args) > 1 {
		DBType = os.Args[1]
	}

	// Choose DB for cats service
	if DBType == "mongo" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Errorf("mongo connect error %v", err)
			return
		}

		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Errorf("mongo ping error %v", err)
			return
		}

		collection := client.Database("local").Collection("cats")
		repoMongo := repository.NewRepoMongo(collection)

		srv = service.NewCatService(repoMongo)

		log.Info("mongo DB is used")
	} else {
		srv = service.NewCatService(repoPostgres)
		log.Info("postgres DB is used")
	}

	cats := handler.NewCatsHandler(srv)

	userRepoPostgres := repository.NewUserRepo(db)
	usersSrv := service.NewUserService(userRepoPostgres)
	users := handler.NewUsersHandler(usersSrv)

	authSrv := service.NewAuthService(userRepoPostgres, cfg)
	auth := handler.NewAuthHandler(authSrv)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, this is a Cats service!")
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
	e.GET("/users/:id", users.GetByID, middleware.JWT(hmacJWTSecret))
	e.GET("/users", users.GetAll, middleware.JWT(hmacJWTSecret))
	e.POST("/users", users.Create, middleware.JWT(hmacJWTSecret))
	e.DELETE("/users/:id", users.Delete, middleware.JWT(hmacJWTSecret))
	e.PUT("/users/:id", users.Update, middleware.JWT(hmacJWTSecret))

	e.Logger.Fatal(e.Start(":1323"))
}
