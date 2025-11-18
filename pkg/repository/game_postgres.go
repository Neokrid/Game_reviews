package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Neokrid/game-review/pkg/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GamePostgres struct {
	db *sqlx.DB
}

func NewGamePostgres(db *sqlx.DB) *GamePostgres {
	return &GamePostgres{db: db}
}

func (r *GamePostgres) CreateGame(game model.Game) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (title, description, developer, release) values ($1, $2, $3, $4) RETURNING id", gamesTable)
	row := r.db.QueryRow(query, game.Title, game.Description, game.Developer, game.Release)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
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
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argsId := 1
	if updateGameArgs.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argsId))
		args = append(args, *updateGameArgs.Title)
		argsId++
	}
	if updateGameArgs.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argsId))
		args = append(args, *updateGameArgs.Description)
		argsId++
	}
	if updateGameArgs.Developer != nil {
		setValues = append(setValues, fmt.Sprintf("developer=$%d", argsId))
		args = append(args, *updateGameArgs.Developer)
		argsId++
	}
	if updateGameArgs.Release != nil {
		releaseDate, err := time.Parse("2006-01-02", *updateGameArgs.Release)
		if err != nil {
			return errors.New("ошибка формата даты(используйте YYYY-MM-DD)")
		}
		setValues = append(setValues, fmt.Sprintf("release=$%d", argsId))
		args = append(args, releaseDate)
		argsId++
	}

	setQuery := strings.Join(setValues, ", ")

	if len(setValues) == 0 {
		return errors.New("структура обновления не имеет полей")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", gamesTable, setQuery, argsId)

	args = append(args, gameId)

	_, err := r.db.Exec(query, args...)
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
