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
		log.Fatalln(err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
		return
	}

	repo := repository.NewRepo(db)

	srv := service.NewCatService(repo)

	cats := handler.NewCat(srv)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, this is a Cats service!")
	})

	e.GET("/cats/:id", cats.GetById)
	e.GET("/cats", cats.GetAll)
	e.POST("/users/:id", cats.Create)
	e.DELETE("/users/:id", cats.Delete)

	e.Logger.Fatal(e.Start(":1323"))
}
