package repository

import (
	"context"
	"database/sql"
	"fmt"
	"golang-postgresql-auth-template/internal/models"
	"log"
)

type UserRepo struct {
	db *sql.DB
	// You can insert a caching service here
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) CreateUser(
	ctx context.Context,
	email string,
	passwordHash string,
	username string,
	firstname string,
	surname string,
) error {
	query := `INSERT INTO users (email, password_hash, username, firstname, surname) VALUES ($1, $2, $3, $4, $5)`
	_, err := u.db.ExecContext(ctx, query, email, passwordHash, username, firstname, surname)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) UpdateUser(ctx context.Context, user models.User) error {
	query := `UPDATE users
			  SET firstname = $1, surname = $2, password_hash = $3
			  WHERE id = $4`

	_, err := u.db.ExecContext(ctx, query,
		user.FirstName, user.SurName,
		user.PasswordHash, user.ID)

	return err
}

func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`
	err := u.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.FirstName,
		&user.SurName,
		&user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1`
	err := u.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.FirstName,
		&user.SurName,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

func (u *UserRepo) CheckIfUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`
	err := u.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		log.Println(err)
	}
	return exists, nil
}

func (u *UserRepo) CheckIfUserExistsByID(ctx context.Context, id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	err := u.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
