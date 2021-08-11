package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Repo struct {
	conn *sql.DB
}

func NewRepo() (*Repo, error) {
	connStr := "user=postgres password=pgpass sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return &Repo{conn: db}, nil
	} else {
		return nil, err
	}
}
