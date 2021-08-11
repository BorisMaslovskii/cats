package protocol

import (
	"github.com/jackc/pgx/v4"
)

type CatsService struct {
	dbConn *pgx.Conn
}

func NewCatsService(dbConn *pgx.Conn) *CatsService {

	return &CatsService{
		dbConn: dbConn,
	}
}
