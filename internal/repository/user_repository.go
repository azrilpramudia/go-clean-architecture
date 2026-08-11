package repository

import (
	"context"
	"database/sql"

	"github.com/azrilpramudia/go-clean-architecture/internal/entity"
)

type UserRespository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
}

type userRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRespository {
	return &userRepositoryImpl{DB: db}
}

func (r *userRepositoryImpl) Save(ctx context.Context, user *entity.User) error {
	query := "INSERT INTO users (username, password, name) VALUES (?, ?, ?)"
	result, err := r.DB.ExecContext(ctx, query, user.Username, user.Password, user.Name)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := "SELECT id, username, password, name FROM users WHERE username = ?"
	row := r.DB.QueryRowContext(ctx, query, username)

	user := new(entity.User)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Name)
	if err != nil {
		return nil, err
	}
	return user, nil
}