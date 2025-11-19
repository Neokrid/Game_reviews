package service

import (
	"github.com/Neokrid/game-review/pkg/model"
	"github.com/Neokrid/game-review/pkg/repository"
	"github.com/google/uuid"
)

type Service struct {
	Authorization
	Game
	Reviews
}

type Authorization interface {
	CreateUser(user model.User) error
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (uuid.UUID, error)
}

type Game interface {
	CreateGame(game model.Game) error
	GetAllGames() ([]model.Game, error)
	GetGamesById(id uuid.UUID) (model.Game, error)
	DeleteGame(gameId uuid.UUID) error
	UpdateGame(gameId uuid.UUID, updateGameArgs model.UpdateGame) error
	GetGamesReviews(gameId uuid.UUID) ([]model.Review, error)
}

type Reviews interface {
	CreateReview(userId, gameId uuid.UUID, input model.Review) error
	GetReviewById(id uuid.UUID) (model.Review, error)
	DeleteReview(id uuid.UUID) error
	UpdateReview(id uuid.UUID, updateReviewArgs model.UpdateReview) error
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Game:          NewGameService(repo.Game),
		Reviews:       NewReviewService(repo.Reviews, repo.Game),
	}
}
