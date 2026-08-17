package repository

import (
	"context"
	"database/sql"

	"github.com/azrilpramudia/go-clean-architecture/internal/entity"
)

type UserRespository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	FindAll(ctx context.Context) ([]entity.User, error)
	Delete(ctx context.Context, id int64) error
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
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	query := "SELECT id, username, password, name FROM users WHERE id = ?"
	row := r.DB.QueryRowContext(ctx, query, id)

	user := new(entity.User)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryImpl) FindAll(ctx context.Context) ([]entity.User, error) {
	query := "SELECT id, username, password, name FROM users"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		user := entity.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}