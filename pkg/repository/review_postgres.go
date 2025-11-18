package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Neokrid/game-review/pkg/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReviewPostgres struct {
	db *sqlx.DB
}

func NewReviewPostgres(db *sqlx.DB) *ReviewPostgres {
	return &ReviewPostgres{db: db}
}

func (r *ReviewPostgres) CreateReview(userId, gameId uuid.UUID, input model.Review) (uuid.UUID, error) {
	var id uuid.UUID
	createReviewQuery := fmt.Sprintf("INSERT INTO %s (game_id, user_id, rating, text_review) VALUES ($1, $2, $3, $4) RETURNING id", reviewTable)
	row := r.db.QueryRow(createReviewQuery, gameId, userId, input.Rating, input.TextReview)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil

}

func (r *ReviewPostgres) GetReviewById(id uuid.UUID) (model.Review, error) {
	var review model.Review
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", reviewTable)
	err := r.db.Get(&review, query, id)
	return review, err
}

func (r *ReviewPostgres) DeleteReview(id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", gamesTable)
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ReviewPostgres) UpdateReview(reviewId uuid.UUID, updateReviewArgs model.UpdateReview) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argsId := 1

	if updateReviewArgs.Rating != nil {
		setValues = append(setValues, fmt.Sprintf("rating=$%d", argsId))
		args = append(args, *updateReviewArgs.Rating)
		argsId++
	}

	if updateReviewArgs.TextReview != nil {
		setValues = append(setValues, fmt.Sprintf("text_review=$%d", argsId))
		args = append(args, *updateReviewArgs.TextReview)
		argsId++
	}

	setQuery := strings.Join(setValues, ", ")

	if len(setValues) == 0 {
		return errors.New("структура обновления не имеет полей")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", reviewTable, setQuery, argsId)

	args = append(args, reviewId)
	_, err := r.db.Exec(query, args...)
	return err
}
