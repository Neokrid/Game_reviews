package model

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" binding:"required" db:"title"`
	Description string    `json:"description" binding:"required" db:"description"`
	Developer   string    `json:"developer" binding:"required" db:"developer"`
	Release     time.Time `json:"release" binding:"required" db:"release"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UpdateGame struct {
	Title       *string `json:"title" db:"title"`
	Description *string `json:"description" db:"description"`
	Developer   *string `json:"developer" db:"developer"`
	Release     *string `json:"release" db:"release"`
}

type GameResponse struct {
	Game       []Game `json:"game"`
	NextCursor string `json:"nextCursor"`
}
