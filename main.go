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
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	var srv *service.CatService

	DBType := ""
	if len(os.Args) > 1 {
		DBType = os.Args[1]
	}

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
		repo := repository.NewRepoMongo(collection)
		srv = service.NewCatService(repo)

		log.Info("mongo DB is used")

	} else {

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

		repo := repository.NewRepo(db)
		srv = service.NewCatService(repo)

		log.Info("postgres DB is used")
	}

	cats := handler.NewCat(srv)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, this is a Cats service!")
	})

	e.GET("/cats/:id", cats.GetByID)
	e.GET("/cats", cats.GetAll)
	e.POST("/cats", cats.Create)
	e.DELETE("/cats/:id", cats.Delete)
	e.PUT("/cats/:id", cats.Update)

	e.Logger.Fatal(e.Start(":1323"))
}
