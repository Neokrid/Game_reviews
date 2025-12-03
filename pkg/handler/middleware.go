package handler

import (
	"net/http"
	"strings"

	e "github.com/Neokrid/game-review/pkg/errors"
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
		e.WriteErr(c, e.NewErr(nil, http.StatusUnauthorized, "Empty authentication header"))
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		e.WriteErr(c, e.NewErr(nil, http.StatusUnauthorized, "Invalid authentication header"))
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		e.WriteErr(c, err)
		return
	}

	c.Set(userCtx, userId)
}

func getUserid(c *gin.Context) (uuid.UUID, error) {
	idValue, ok := c.Get(userCtx)
	if !ok {
		return uuid.Nil, e.NewErr(nil, http.StatusUnauthorized, "user ID not found in context")
	}
	userID, ok := idValue.(uuid.UUID)
	if !ok {
		return uuid.Nil, e.NewErr(nil, http.StatusUnauthorized, "user ID in context is of invalid type")
	}

	return userID, nil
}
