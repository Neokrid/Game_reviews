package repository

import (
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Neokrid/game-review/pkg/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReviewPostgres struct {
	db *sqlx.DB
}

func NewReviewPostgres(db *sqlx.DB) *ReviewPostgres {
	return &ReviewPostgres{
		db: db,
	}
}

func (r *ReviewPostgres) CreateReview(userId, gameId uuid.UUID, input model.Review) error {
	var id uuid.UUID
	createReviewQuery := fmt.Sprintf("INSERT INTO %s (game_id, user_id, rating, text_review) VALUES ($1, $2, $3, $4) RETURNING id", reviewTable)
	row := r.db.QueryRow(createReviewQuery, gameId, userId, input.Rating, input.TextReview)
	if err := row.Scan(&id); err != nil {
		return err
	}

	return nil

}

func (r *ReviewPostgres) GetReviewById(id uuid.UUID) (model.Review, error) {
	var review model.Review
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", reviewTable)
	err := r.db.Get(&review, query, id)
	return review, err
}

func (r *ReviewPostgres) DeleteReview(id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", reviewTable)
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ReviewPostgres) UpdateReview(reviewId uuid.UUID, updateReviewArgs model.UpdateReview) error {
	query := sq.Update("reviews").Where(sq.Eq{"id": reviewId}).PlaceholderFormat(sq.Dollar)
	isUpdate := false

	if updateReviewArgs.TextReview != nil {
		query = query.Set("text_review", *updateReviewArgs.TextReview)
		isUpdate = true
	}

	if updateReviewArgs.Rating != nil {
		query = query.Set("rating", *updateReviewArgs.Rating)
		isUpdate = true
	}

	if !isUpdate {
		return errors.New("структура обновления не имеет полей")
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(sqlQuery, args...)
	return err
}
