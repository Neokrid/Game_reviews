package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" binding:"required" db:"name"`
	UserName     string    `json:"username" binding:"required" db:"username"`
	PasswordHash string    `json:"password_hash" binding:"required" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}
