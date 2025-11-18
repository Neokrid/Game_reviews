package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "пустой заголовок аутентификации")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "неверный заголовок аутентификации")
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

func getUserid(c *gin.Context) (uuid.UUID, error) {
	idValue, ok := c.Get(userCtx)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	userID, ok := idValue.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID in context is of invalid type")
	}

	return userID, nil
}
