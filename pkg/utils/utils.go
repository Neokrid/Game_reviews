package utils

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Neokrid/game-review/pkg/errors"
	"github.com/Neokrid/game-review/pkg/model"
	"github.com/google/uuid"
)

func EncodeCursor(id uuid.UUID) (string, error) {
	c := model.Cursor{
		GameId: id,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func DecodeCursor(token string) (uuid.UUID, error) {
	if token == "" {
		return uuid.Nil, nil
	}
	bytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return uuid.Nil, errors.NewErr(nil, http.StatusBadRequest, "Invalid token format")
	}
	var c model.Cursor
	if err := json.Unmarshal(bytes, &c); err != nil {
		return uuid.Nil, err
	}
	return c.GameId, nil
}
