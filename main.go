package main

import (
	"database/sql"
	"net/http"

	"github.com/BorisMaslovskii/cats/internal/handler"
	"github.com/BorisMaslovskii/cats/internal/repository"
	"github.com/BorisMaslovskii/cats/internal/service"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	connStr := "user=postgres password=pgpass sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Errorf("sql Open %v", err)
		return
	}
	defer func() {
		errd := db.Close()
		if err != nil {
			log.Errorf("defer db Close %v", errd)
		}
	}()

	err = db.Ping()
	if err != nil {
		log.Errorf("db Ping %v", err)
		return
	}

	repo := repository.NewRepo(db)

	srv := service.NewCatService(repo)

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
