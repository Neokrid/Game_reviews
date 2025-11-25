package repository

import (
	"time"

	"github.com/go-redis/redis"
)

type GetGameRedis struct {
	redis *redis.Client
}

func NewGameRedis(redis *redis.Client) *GetGameRedis {
	return &GetGameRedis{
		redis: redis,
	}
}

func (r *GetGameRedis) GetLeaderboardCache(cacheKay string) (string, error) {
	return r.redis.Get(cacheKay).Result()

}

func (r *GetGameRedis) SetLeaderboardCache(cacheKay string, data []byte) error {
	err := r.redis.Set(cacheKay, data, 24*time.Hour)
	if err != nil {
		return err.Err()
	}
	return redis.Nil
}
