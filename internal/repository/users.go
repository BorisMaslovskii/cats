package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BorisMaslovskii/cats/internal/model"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// UserRepository interface provides crud functions for the user service not depending on the database type
type UserRepository interface {
	GetAll(ctx context.Context) ([]*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Create(ctx context.Context, user *model.User) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, user *model.User) error
	GetByLogin(ctx context.Context, reqUser *model.User) (*model.User, error)
}

// userRepositoryPostgres struct
type userRepositoryPostgres struct {
	conn *sql.DB
}

// NewUserRepo func creates new UserRepository
func NewUserRepo(conn *sql.DB) UserRepository {
	return &userRepositoryPostgres{conn: conn}
}

// GetAll gets all users from userRepositoryPostgres
func (r *userRepositoryPostgres) GetAll(ctx context.Context) ([]*model.User, error) {
	rows, err := r.conn.QueryContext(ctx, "select id, login, password, admin from users")
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Errorf("users GetAll rows Close %v", err)
		}
	}()
	users := []*model.User{}

	for rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.ID, &user.Login, &user.Password, &user.Admin)
		if err != nil {
			log.Errorf("users GetAll rows.Scan %v", err)
			continue
		}
		users = append(users, &user)
	}

	return users, nil
}

// GetByID func gets a user by id from userRepositoryPostgres
func (r *userRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row := r.conn.QueryRowContext(ctx, "select id, login, password, admin from users where id = $1", id)
	user := &model.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Admin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// Create func creates a new user into userRepositoryPostgres
func (r *userRepositoryPostgres) Create(ctx context.Context, user *model.User) (uuid.UUID, error) {
	id := uuid.New()

	_, err := r.conn.ExecContext(ctx, "insert into users(id, login, password, admin) values ($1, $2, $3, $4)", id, user.Login, user.Password, user.Admin)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// Delete func deletes a user from userRepositoryPostgres
func (r *userRepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.conn.ExecContext(ctx, "delete from users where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// Update func updates a user in userRepositoryPostgres
func (r *userRepositoryPostgres) Update(ctx context.Context, id uuid.UUID, user *model.User) error {
	_, err := r.conn.ExecContext(ctx, "update users set login = $1, password = $2, admin = $3 where id = $4", user.Login, user.Password, user.Admin, id)
	if err != nil {
		return err
	}
	return nil
}

// Update func updates a user in userRepositoryPostgres
func (r *userRepositoryPostgres) GetByLogin(ctx context.Context, reqUser *model.User) (*model.User, error) {
	row := r.conn.QueryRowContext(ctx, "select id, login, password, admin from users where login = $1", reqUser.Login)
	user := &model.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Admin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return user, nil
}
