package repository

import (
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
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

func (r *GamePostgres) GetAllGames() ([]model.Game, error) {
	var games []model.Game
	query := fmt.Sprintf("SELECT * FROM %s", gamesTable)
	err := r.db.Select(&games, query)

	return games, err
}

func (r *GamePostgres) GetGamesById(id uuid.UUID) (model.Game, error) {
	var game model.Game
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", gamesTable)
	err := r.db.Get(&game, query, id)
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
			return errors.New("ошибка формата даты(используйте YYYY-MM-DD)")
		}
		query = query.Set("release", releaseDate)
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
	query := fmt.Sprintf("SELECT ROW_NUMBER() OVER (ORDER BY CAST(AVG(r.rating) AS DECIMAL(10, 2)) DESC, COUNT(r.id) DESC) as position, g.title, CAST(AVG(r.rating) AS DECIMAL(10, 2)) AS average_rating from %s g join %s r on g.id = r.game_id group by g.id, g.title", gamesTable, reviewTable)
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
	query := fmt.Sprintf("SELECT created_at::DATE AS review_date, ROUND(AVG(rating), 1) AS avg_rating FROM %s WHERE game_id = $1 GROUP BY review_date ORDER BY review_date ASC ", reviewTable)
	if err := r.db.Select(&ratingHistory, query, gameId); err != nil {
		return nil, err
	}
	return ratingHistory, nil
}
