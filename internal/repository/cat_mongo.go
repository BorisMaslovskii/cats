package repository

import "database/sql"

// catRepository struct
type catRepositoryMongo struct {
	conn *sql.DB
}
