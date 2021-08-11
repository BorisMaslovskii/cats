package main

import (
	"database/sql"

	"github.com/BorisMaslovskii/cats/internal/handler"
	"github.com/BorisMaslovskii/cats/internal/repository"
	"github.com/BorisMaslovskii/cats/internal/service"
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

	repo := repository.NewRepo(db)

	srv := service.NewCatService(repo)

	h := handler.NewCat(srv)
}
