package repository

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Neokrid/game-review/pkg/errors"
	"github.com/Neokrid/game-review/pkg/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GamePostgres struct {
	db *sqlx.DB
}

func NewGamePostgres(db *sqlx.DB) *GamePostgres {
	return &GamePostgres{
		db: db,
	}
}

func (r *GamePostgres) CreateGame(game model.Game) error {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (title, description, developer, release) values ($1, $2, $3, $4) RETURNING id", gamesTable)
	row := r.db.QueryRow(query, game.Title, game.Description, game.Developer, game.Release)
	if err := row.Scan(&id); err != nil {
		return err
	}
	return nil
}

func (r *GamePostgres) GetAllGames(limit int, lastId uuid.UUID) ([]model.Game, error) {
	var games []model.Game
	query := sq.Select("*").From("games").PlaceholderFormat(sq.Dollar)
	if lastId != uuid.Nil {
		query = query.Where(sq.Lt{"id": lastId})
	}
	query = query.OrderBy("id DESC")
	query = query.Limit(uint64(limit))
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	if err := r.db.Select(&games, sqlQuery, args...); err != nil {
		return nil, err
	}
	return games, nil
}

func (r *GamePostgres) GetGamesById(id uuid.UUID) (model.Game, error) {
	var game model.Game
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", gamesTable)
	err := r.db.Get(&game, query, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return game, errors.NewErr(err, http.StatusNotFound, "The requested game was not found.")
		}
	}
	return game, err
}

func (r *GamePostgres) DeleteGame(gameId uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", gamesTable)
	_, err := r.db.Exec(query, gameId)
	return err

}

func (r *GamePostgres) UpdateGame(gameId uuid.UUID, updateGameArgs model.UpdateGame) error {
	query := sq.Update("games").Where(sq.Eq{"id": gameId}).PlaceholderFormat(sq.Dollar)
	isUpdate := false
	if updateGameArgs.Title != nil {
		query = query.Set("title", *updateGameArgs.Title)
		isUpdate = true
	}
	if updateGameArgs.Description != nil {
		query = query.Set("description", *updateGameArgs.Description)
		isUpdate = true
	}
	if updateGameArgs.Developer != nil {
		query = query.Set("developer", *updateGameArgs.Developer)
		isUpdate = true
	}
	if updateGameArgs.Release != nil {
		releaseDate, err := time.Parse("2006-01-02", *updateGameArgs.Release)
		if err != nil {
			return errors.NewErr(nil, http.StatusBadRequest, "Invalid Date format. Format: 2006-01-02")
		}
		query = query.Set("release", releaseDate)
		isUpdate = true
	}

	if !isUpdate {
		return errors.NewErr(nil, http.StatusBadRequest, "The structure has no updated fields.")
	}
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(sqlQuery, args...)
	return err
}

func (r *GamePostgres) GetGamesReviews(gameId uuid.UUID) ([]model.Review, error) {
	reviews := make([]model.Review, 0)
	query := fmt.Sprintf("SELECT id, game_id, user_id, rating, text_review, created_at FROM %s WHERE game_id = $1", reviewTable)
	if err := r.db.Select(&reviews, query, gameId); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *GamePostgres) GetLeaderboard() ([]model.Leaderboard, error) {
	leaderboard := make([]model.Leaderboard, 0)
	query := fmt.Sprintf(`
	SELECT 
	ROW_NUMBER() OVER (ORDER BY CAST(AVG(r.rating) AS DECIMAL(10, 2)) DESC, 
	COUNT(r.id) DESC) as position, g.title, 
	CAST(AVG(r.rating) AS DECIMAL(10, 2)) AS average_rating 
	from %s g join %s r on g.id = r.game_id group by g.id, g.title LIMIT 10`,
		gamesTable, reviewTable)
	err := r.db.Select(&leaderboard, query)
	return leaderboard, err
}

func (r *GamePostgres) SearchGame(gameToFind model.Game) ([]model.Game, error) {
	gamesFound := make([]model.Game, 0)
	query := fmt.Sprintf("SELECT id, title, description, developer from %s where title %% $1 order by title ASC limit 10", gamesTable)
	if err := r.db.Select(&gamesFound, query, gameToFind.Title); err != nil {
		return nil, err
	}

	return gamesFound, nil
}

func (r *GamePostgres) GetRatingHistory(gameId uuid.UUID) ([]model.RatingHistory, error) {
	ratingHistory := make([]model.RatingHistory, 0)
	query := fmt.Sprintf(`
	SELECT 
	created_at::DATE AS review_date, 
	ROUND(AVG(rating), 1) AS avg_rating 
	FROM %s WHERE game_id = $1 
	GROUP BY review_date 
	ORDER BY review_date ASC`, reviewTable)
	if err := r.db.Select(&ratingHistory, query, gameId); err != nil {
		return nil, err
	}

	return ratingHistory, nil
}
