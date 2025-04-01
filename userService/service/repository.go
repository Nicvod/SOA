package main

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (int, error)
	GetUserByLogin(ctx context.Context, login string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id int) (*User, error)
}

type UserRepositorySpec struct {
	db *sqlx.DB
}

type User struct {
	ID          int       `json:"id" db:"id"`
	Login       string    `json:"login" db:"login"`
	Password    string    `json:"-" db:"password"`
	Email       string    `json:"email" db:"email"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	BirthDate   time.Time `json:"birth_date" db:"birth_date"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &UserRepositorySpec{db: db}
}

func (r *UserRepositorySpec) CreateUser(ctx context.Context, user *User) (int, error) {
	query := `
        INSERT INTO users (login, password, email, first_name, last_name, birth_date, phone_number, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `
	var id int
	err := r.db.QueryRowContext(ctx, query,
		user.Login, user.Password, user.Email, user.FirstName, user.LastName, user.BirthDate, user.PhoneNumber, user.CreatedAt, user.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepositorySpec) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	var user User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE login = $1", login)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositorySpec) UpdateUser(ctx context.Context, user *User) error {
	query := `
        UPDATE users
        SET email = $1, first_name = $2, last_name = $3, birth_date = $4, phone_number = $5, updated_at = $6
        WHERE id = $7
    `
	_, err := r.db.ExecContext(ctx, query,
		user.Email, user.FirstName, user.LastName, user.BirthDate, user.PhoneNumber, user.UpdatedAt, user.ID,
	)
	return err
}

func (r *UserRepositorySpec) GetUserByID(ctx context.Context, id int) (*User, error) {
	var user User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
