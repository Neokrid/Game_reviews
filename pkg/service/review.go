package service

import (
	"github.com/Neokrid/game-review/pkg/model"
	"github.com/Neokrid/game-review/pkg/repository"
	"github.com/google/uuid"
)

type ReviewService struct {
	repo     repository.Reviews
	gameRepo repository.Game
}

func NewReviewService(repo repository.Reviews, gameRepo repository.Game) *ReviewService {
	return &ReviewService{
		repo:     repo,
		gameRepo: gameRepo,
	}
}

func (s *ReviewService) CreateReview(userId, gameId uuid.UUID, input model.Review) (uuid.UUID, error) {
	_, err := s.gameRepo.GetGamesById(gameId)
	if err != nil {
		return uuid.Nil, err
	}
	return s.repo.CreateReview(userId, gameId, input)
}

func (s *ReviewService) GetReviewById(id uuid.UUID) (model.Review, error) {
	return s.repo.GetReviewById(id)
}

func (s *ReviewService) DeleteReview(id uuid.UUID) error {
	return s.repo.DeleteReview(id)
}

func (s *ReviewService) UpdateReview(id uuid.UUID, updateReviewArgs model.UpdateReview) error {
	return s.repo.UpdateReview(id, updateReviewArgs)
}
