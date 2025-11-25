package repository

import "github.com/go-redis/redis"

type ConfigRedis struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisDB(cfg ConfigRedis) (*redis.Client, error) {
	redis := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := redis.Ping().Err(); err != nil {
		return nil, err
	}
	return redis, nil
}
