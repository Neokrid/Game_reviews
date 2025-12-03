package handler

import (
	"net/http"
	"time"

	"github.com/Neokrid/game-review/pkg/errors"
	"github.com/Neokrid/game-review/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getAllGamesResponse struct {
	Games []model.Game `json:"Games"`
}

type GetLeaderboardResponse struct {
	Leaderboard []model.Leaderboard `json:"leaderboard"`
}

func (h *Handler) getAllGames(c *gin.Context) {
	cursorToken := c.Query("token")
	limitStr := c.DefaultQuery("limit", "10")
	games, err := h.services.Game.GetAllGames(limitStr, cursorToken)
	if err != nil {
		errors.WriteErr(c, err)
		return
	}
	c.JSON(http.StatusOK, games)
}

func (h *Handler) getGamesById(c *gin.Context) {
	id := c.Param("id")
	gameId, err := uuid.Parse(id)
	if err != nil {
		errors.WriteErr(c, errors.NewErr(nil, http.StatusUnauthorized, "Invalid ID format"))
		return
	}

	game, err := h.services.Game.GetGamesById(gameId)
	if err != nil {
		errors.WriteErr(c, err)
		return
	}
	c.JSON(http.StatusOK, game)
}

type GameReviewsResponse struct {
	GameTitle string         `json:"game_title"`
	Reviews   []model.Review `json:"reviews"`
}

func (h *Handler) getGamesReviews(c *gin.Context) {
	gameIdStr := c.Param("id")
	gameId, err := uuid.Parse(gameIdStr)
	if err != nil {
		errors.WriteErr(c, errors.NewErr(nil, http.StatusBadRequest, "Invalid ID format"))
		return
	}

	game, err := h.services.Game.GetGamesById(gameId)
	if err != nil {
		errors.WriteErr(c, err)
		return
	}

	reviews, err := h.services.Game.GetGamesReviews(gameId)
	if err != nil {
		errors.WriteErr(c, err)
		return
	}

	response := GameReviewsResponse{
		GameTitle: game.Title,
		Reviews:   reviews,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) deleteGame(c *gin.Context) {
	id := c.Param("id")
	gameId, err := uuid.Parse(id)
	if err != nil {
		errors.WriteErr(c, errors.NewErr(nil, http.StatusBadRequest, "Invalid ID format"))
		return
	}

	err = h.services.Game.DeleteGame(gameId)
	if err != nil {
		errors.WriteErr(c, err)
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

type createGameInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Developer   string `json:"developer" binding:"required"`
	Release     string `json:"release" binding:"required"`
}

func (h *Handler) createGame(c *gin.Context) {
	var input createGameInput

	if err := c.BindJSON(&input); err != nil {
		errors.WriteErr(c, err)
		return
	}

	releaseDate, err := time.Parse("2006-01-02", input.Release)
	if err != nil {
		errors.WriteErr(c, errors.NewErr(nil, http.StatusBadRequest, "Invalid Date format. Format: 2006-01-02"))
		return
	}

	gameArg := model.Game{
		Title:       input.Title,
		Description: input.Description,
		Developer:   input.Developer,
		Release:     releaseDate,
	}

	err = h.services.Game.CreateGame(gameArg)
	if err != nil {
		errors.WriteErr(c, err)
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) changeGame(c *gin.Context) {
	id := c.Param("id")
	gameId, err := uuid.Parse(id)
	if err != nil {
		errors.WriteErr(c, errors.NewErr(nil, http.StatusBadRequest, "Invalid ID format"))
		return
	}
	var input model.UpdateGame
	if err := c.BindJSON(&input); err != nil {
		errors.WriteErr(c, err)
		return
	}

	if err := h.services.UpdateGame(gameId, input); err != nil {
		errors.WriteErr(c, err)
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) getLeaderboard(c *gin.Context) {
	leaderboard, err := h.services.GetLeaderboard()
	if err != nil {
		errors.WriteErr(c, err)
		return
	}
	c.JSON(http.StatusOK, GetLeaderboardResponse{
		Leaderboard: leaderboard,
	})

}

type SearchGameInput struct {
	Title string `json:"title" binding:"required"`
}

func (h *Handler) searchGame(c *gin.Context) {
	var input SearchGameInput
	if err := c.BindJSON(&input); err != nil {
		errors.WriteErr(c, err)
		return
	}

	gameToFind := model.Game{
		Title: input.Title,
	}
	gamesFound, err := h.services.SearchGame(gameToFind)
	if err != nil {
		errors.WriteErr(c, err)
		return
	}
	c.JSON(http.StatusOK, getAllGamesResponse{
		Games: gamesFound,
	})
}

func (h *Handler) getRatingHistory(c *gin.Context) {
	id := c.Param("id")
	gameId, err := uuid.Parse(id)
	if err != nil {
		errors.WriteErr(c, errors.NewErr(nil, http.StatusBadRequest, "Invalid ID format"))
		return
	}

	ratingHistory, err := h.services.GetRatingHistory(gameId)
	if err != nil {
		errors.WriteErr(c, err)
		return
	}
	c.JSON(http.StatusOK, ratingHistory)
}
