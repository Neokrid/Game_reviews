package repository

import (
	"encoding/json"
	"time"

	"github.com/Neokrid/game-review/pkg/model"
	"github.com/go-redis/redis"
)

const (
	CacheTTL = 24 * time.Hour
)

type GetGameRedis struct {
	redis *redis.Client
}

func NewGameRedis(redis *redis.Client) *GetGameRedis {
	return &GetGameRedis{
		redis: redis,
	}
}

func (r *GetGameRedis) GetLeaderboardCache(cacheKay string) ([]model.Leaderboard, error) {
	val, err := r.redis.Get(cacheKay).Bytes()
	if err != nil {
		return nil, err
	}
	var leaderboard []model.Leaderboard
	if err := json.Unmarshal(val, &leaderboard); err != nil {
		return nil, err
	}

	return leaderboard, nil

}

func (r *GetGameRedis) SetLeaderboardCache(cacheKay string, leaderboard []model.Leaderboard) error {
	data, err := json.Marshal(leaderboard)
	if err != nil {
		return err
	}
	errRedis := r.redis.Set(cacheKay, data, CacheTTL)
	if errRedis != nil {
		return errRedis.Err()
	}
	return nil
}
