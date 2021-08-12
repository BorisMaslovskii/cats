package repository

import (
	"context"
	"database/sql"

	"github.com/BorisMaslovskii/cats/internal/model"
	log "github.com/sirupsen/logrus"
)

type Repository struct {
	conn *sql.DB
}

func NewRepo(conn *sql.DB) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) GetAllCats(ctx context.Context) ([]*model.Cat, error) {
	rows, err := r.conn.QueryContext(ctx, "select ID, NAME, COLOR from CATS")
	if err != nil {
		return []*model.Cat{}, err
	}
	defer rows.Close()
	cats := []*model.Cat{}

	for rows.Next() {
		cat := model.Cat{}
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Color)
		if err != nil {
			log.Errorf("GetAllCats rows.Scan %v", err)
			continue
		}
		cats = append(cats, &cat)
	}

	return cats, nil
}
