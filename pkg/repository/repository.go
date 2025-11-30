package repository

import (
	"github.com/Neokrid/game-review/pkg/model"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Authorization
	Game
	Reviews
	GameRedis
}

type Authorization interface {
	CreateUser(user model.User) error
	GetUser(username string) (model.User, error)
}

type Game interface {
	CreateGame(game model.Game) error
	GetAllGames() ([]model.Game, error)
	GetGamesById(id uuid.UUID) (model.Game, error)
	DeleteGame(gameId uuid.UUID) error
	UpdateGame(gameId uuid.UUID, updateGameArgs model.UpdateGame) error
	GetGamesReviews(gameId uuid.UUID) ([]model.Review, error)
	GetLeaderboard() ([]model.Leaderboard, error)
	SearchGame(gameToFind model.Game) ([]model.Game, error)
	GetRatingHistory(gameId uuid.UUID) ([]model.RatingHistory, error)
}

type Reviews interface {
	CreateReview(userId, gameId uuid.UUID, input model.Review) error
	GetReviewById(id uuid.UUID) (model.Review, error)
	DeleteReview(id uuid.UUID) error
	UpdateReview(id uuid.UUID, updateReviewArgs model.UpdateReview) error
}

type GameRedis interface {
	GetLeaderboardCache(cacheKay string) ([]model.Leaderboard, error)
	SetLeaderboardCache(cacheKay string, leaderboard []model.Leaderboard) error
}

func NewRepository(db *sqlx.DB, redis *redis.Client) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Game:          NewGamePostgres(db),
		Reviews:       NewReviewPostgres(db),
		GameRedis:     NewGameRedis(redis),
	}
}
