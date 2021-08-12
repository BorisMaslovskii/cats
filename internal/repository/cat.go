// Package repository provides repository (database, filesystem etc) logic
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BorisMaslovskii/cats/internal/model"
	log "github.com/sirupsen/logrus"
)

type CatRepository interface {
	GetAll(ctx context.Context) ([]*model.Cat, error)
	GetByID(ctx context.Context, id string) (*model.Cat, error)
	Create(ctx context.Context, cat *model.Cat) (int, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, cat *model.Cat) error
}

// catRepository struct
type catRepository struct {
	conn *sql.DB
}

// NewRepo func creates new catRepository
func NewRepo(conn *sql.DB) *catRepository {
	return &catRepository{conn: conn}
}

// GetAll gets all cats from catRepository
func (r *catRepository) GetAll(ctx context.Context) ([]*model.Cat, error) {
	rows, err := r.conn.QueryContext(ctx, "select id, name, color from cats")
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Errorf("rows Close %v", err)
		}
	}()
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

// GetByID func gets a cat by id from catRepository
func (r *catRepository) GetByID(ctx context.Context, id string) (*model.Cat, error) {
	row := r.conn.QueryRowContext(ctx, "select id, name, color from cats where id = $1", id)
	cat := &model.Cat{}
	err := row.Scan(&cat.ID, &cat.Name, &cat.Color)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return cat, nil
}

// Create func creates a new cat into catRepository
func (r *catRepository) Create(ctx context.Context, cat *model.Cat) (int, error) {
	catID := 0
	err := r.conn.QueryRowContext(ctx, "insert into cats(name, color) values ($1, $2) returning id;", cat.Name, cat.Color).Scan(&catID)
	if err != nil {
		return 0, err
	}
	return catID, nil
}

// Delete func deletes a cat from catRepository
func (r *catRepository) Delete(ctx context.Context, id string) error {
	_, err := r.conn.ExecContext(ctx, "delete from cats where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// Update func updates a cat in catRepository
func (r *catRepository) Update(ctx context.Context, id string, cat *model.Cat) error {
	_, err := r.conn.ExecContext(ctx, "update cats set name = $1, color = $2 where id = $3", cat.Name, cat.Color, id)
	if err != nil {
		return err
	}
	return nil
}
