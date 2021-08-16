package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

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

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.
var hmacSampleSecret []byte

func main() {
	var srv *service.CatService

	connStr := "user=postgres password=pgpass sslmode=disable"
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

	cats := handler.NewCat(srv)

	userRepoPostgres := repository.NewUserRepo(db)
	usersSrv := service.NewUserService(userRepoPostgres)
	users := handler.NewUser(usersSrv)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, this is a Cats service!")
	})

	// JWT
	hmacSampleSecret = []byte("testSecret")

	// Cats service
	e.GET("/cats/:id", cats.GetByID, middleware.JWT(hmacSampleSecret))
	e.GET("/cats", cats.GetAll, middleware.JWT(hmacSampleSecret))
	e.POST("/cats", cats.Create, middleware.JWT(hmacSampleSecret))
	e.DELETE("/cats/:id", cats.Delete, middleware.JWT(hmacSampleSecret))
	e.PUT("/cats/:id", cats.Update, middleware.JWT(hmacSampleSecret))

	// Users service
	e.GET("/users/:id", users.GetByID, middleware.JWT(hmacSampleSecret))
	e.GET("/users", users.GetAll, middleware.JWT(hmacSampleSecret))
	e.POST("/users", users.Create, middleware.JWT(hmacSampleSecret))
	e.DELETE("/users/:id", users.Delete, middleware.JWT(hmacSampleSecret))
	e.PUT("/users/:id", users.Update, middleware.JWT(hmacSampleSecret))
	e.GET("/users/login", users.LogIn)
	e.POST("/users/login", users.LogIn)

	e.Logger.Fatal(e.Start(":1323"))
}
