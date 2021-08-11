package repository

import (
	"database/sql"
)

type Repo struct {
	conn *sql.DB
}

func NewRepo(conn *sql.DB) *Repo {
	return &Repo{conn: conn}
}
