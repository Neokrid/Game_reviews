package model

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	Id         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	GameID     uuid.UUID `json:"game_id" db:"game_id"`
	Rating     int       `json:"rating" binding:"required" db:"rating"`
	TextReview string    `json:"text_review" binding:"required" db:"text_review"`
	CreatedAt  time.Time `json:"create_at" db:"created_at"`
}

type UpdateReview struct {
	Rating     *int    `json:"rating" db:"rating"`
	TextReview *string `json:"text_review" db:"text_review"`
}
