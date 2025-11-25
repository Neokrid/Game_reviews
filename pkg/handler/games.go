package handler

import (
	"net/http"
	"time"

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
	games, err := h.services.Game.GetAllGames()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllGamesResponse{
		Games: games,
	})
}

func (h *Handler) getGamesById(c *gin.Context) {
	id := c.Param("id")
	gameId, err := uuid.Parse(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "неверный формат Id")
	}

	game, err := h.services.Game.GetGamesById(gameId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusBadRequest, "неверный формат Id")
		return
	}

	game, err := h.services.Game.GetGamesById(gameId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	reviews, err := h.services.Game.GetGamesReviews(gameId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusBadRequest, "неверный формат Id")
	}

	err = h.services.Game.DeleteGame(gameId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	releaseDate, err := time.Parse("2006-01-02", input.Release)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
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
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusBadRequest, "неверный формат Id")
		return
	}
	var input model.UpdateGame
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.UpdateGame(gameId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) getLeaderboard(c *gin.Context) {
	leaderboard, err := h.services.GetLeaderboard()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	gameToFind := model.Game{
		Title: input.Title,
	}
	gamesFound, err := h.services.SearchGame(gameToFind)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllGamesResponse{
		Games: gamesFound,
	})
}
