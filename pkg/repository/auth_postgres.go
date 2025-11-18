package repository

import (
	"fmt"

	"github.com/Neokrid/game-review/pkg/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user model.User) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1,$2,$3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.UserName, user.PasswordHash)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id, password_hash  FROM %s WHERE username = $1", usersTable)
	err := r.db.Get(&user, query, username)
	return user, err
}
