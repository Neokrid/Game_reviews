package handler

import (
	"net/http"

	"github.com/Neokrid/game-review/pkg/errors"
	"github.com/Neokrid/game-review/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) getReviewById(c *gin.Context) {
	id := c.Param("id")
	reviewId, err := uuid.Parse(id)
	if err != nil {
		errors.WriteErr(c, errors.NewErr(nil, http.StatusBadRequest, "Invalid ID format"))
		return
	}

	review, err := h.services.Reviews.GetReviewById(reviewId)
	if err != nil {
		errors.WriteErr(c, err)
		return
	}
	c.JSON(http.StatusOK, review)
}

func (h *Handler) deleteReview(c *gin.Context) {
	id := c.Param("review_id")
	reviewID, err := uuid.Parse(id)
	if err != nil {
		errors.WriteErr(c, errors.NewErr(nil, http.StatusBadRequest, "Invalid ID format"))
		return
	}

	err = h.services.Reviews.DeleteReview(reviewID)
	if err != nil {
		errors.WriteErr(c, err)
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) createReview(c *gin.Context) {
	userId, err := getUserid(c)
	if err != nil {
		errors.WriteErr(c, err)
		return
	}

	gameIdStr := c.Param("id")
	gameId, err := uuid.Parse(gameIdStr)
	if err != nil {
		errors.WriteErr(c, errors.NewErr(nil, http.StatusBadRequest, "Invalid ID format"))
		return
	}

	var input model.Review
	if err := c.BindJSON(&input); err != nil {
		errors.WriteErr(c, err)
		return
	}

	err = h.services.Reviews.CreateReview(userId, gameId, input)
	if err != nil {
		errors.WriteErr(c, err)
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) changeReview(c *gin.Context) {
	id := c.Param("review_id")
	reviewId, err := uuid.Parse(id)
	if err != nil {
		errors.WriteErr(c, errors.NewErr(nil, http.StatusBadRequest, "Invalid ID format"))
		return
	}
	var input model.UpdateReview
	if err := c.BindJSON(&input); err != nil {
		errors.WriteErr(c, err)
		return
	}

	if err := h.services.UpdateReview(reviewId, input); err != nil {
		errors.WriteErr(c, err)
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
