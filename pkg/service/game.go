package service

import (
	"strconv"

	"github.com/Neokrid/game-review/pkg/model"
	"github.com/Neokrid/game-review/pkg/repository"
	"github.com/Neokrid/game-review/pkg/utils"
	"github.com/google/uuid"
)

type GameService struct {
	repo      repository.Game
	repoRedis repository.GameRedis
}

func NewGameService(repo repository.Game, repoRedis repository.GameRedis) *GameService {
	return &GameService{repo: repo, repoRedis: repoRedis}
}

func (s *GameService) CreateGame(game model.Game) error {
	return s.repo.CreateGame(game)
}

func (s *GameService) GetAllGames(limitStr string, token string) (*model.GameResponse, error) {
	limit, _ := strconv.Atoi(limitStr)
	lastId, err := utils.DecodeCursor(token)
	if err != nil {
		return nil, err
	}

	if limit < 1 || limit > 20 {
		limit = 20
	}

	games, err := s.repo.GetAllGames(limit, lastId)
	if err != nil {
		return nil, err
	}
	res := &model.GameResponse{
		Game: games,
	}

	if len(games) == limit {
		lastGame := games[len(games)-1]
		res.NextCursor, _ = utils.EncodeCursor(lastGame.Id)
	}
	return res, nil
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

func (s *GameService) GetLeaderboard() ([]model.Leaderboard, error) {
	cacheKey := "leaderboard:"
	val, err := s.repoRedis.GetLeaderboardCache(cacheKey)
	if err == nil {
		return val, nil
	}

	leaderboard, err := s.repo.GetLeaderboard()
	if err != nil {
		return nil, err
	}
	go func() {
		s.repoRedis.SetLeaderboardCache(cacheKey, leaderboard)
	}()
	return leaderboard, nil
}

func (s *GameService) SearchGame(gameToFind model.Game) ([]model.Game, error) {
	return s.repo.SearchGame(gameToFind)
}

func (s *GameService) GetRatingHistory(gameId uuid.UUID) ([]model.RatingHistory, error) {
	return s.repo.GetRatingHistory(gameId)
}
