package handler

import (
	"net/http"

	"github.com/Neokrid/game-review/pkg/errors"
	"github.com/Neokrid/game-review/pkg/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input model.User

	if err := c.ShouldBindJSON(&input); err != nil {
		errors.WriteErr(c, err)
		return
	}
	err := h.services.Authorization.CreateUser(input)
	if err != nil {
		errors.WriteErr(c, err)
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})

}

type SingInInput struct {
	UserName     string `json:"username" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input SingInInput

	if err := c.BindJSON(&input); err != nil {
		errors.WriteErr(c, err)
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.UserName, input.PasswordHash)
	if err != nil {
		errors.WriteErr(c, errors.NewErr(nil, http.StatusUnauthorized, "Incorrect login or password"))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})

}
