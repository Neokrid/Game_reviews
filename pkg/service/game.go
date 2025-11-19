package service

import (
	"github.com/Neokrid/game-review/pkg/model"
	"github.com/Neokrid/game-review/pkg/repository"
	"github.com/google/uuid"
)

type GameService struct {
	repo repository.Game
}

func NewGameService(repo repository.Game) *GameService {
	return &GameService{repo: repo}
}

func (s *GameService) CreateGame(game model.Game) error {
	return s.repo.CreateGame(game)
}

func (s *GameService) GetAllGames() ([]model.Game, error) {
	return s.repo.GetAllGames()
}

func (s *GameService) GetGamesById(id uuid.UUID) (model.Game, error) {
	return s.repo.GetGamesById(id)
}

func (s *GameService) DeleteGame(gameId uuid.UUID) error {
	return s.repo.DeleteGame(gameId)
}

func (s *GameService) UpdateGame(gameId uuid.UUID, updateGameArgs model.UpdateGame) error {
	return s.repo.UpdateGame(gameId, updateGameArgs)
}

func (s *GameService) GetGamesReviews(gameId uuid.UUID) ([]model.Review, error) {
	return s.repo.GetGamesReviews(gameId)
}
