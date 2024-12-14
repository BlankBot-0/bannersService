package cache

import (
	"banners/internal/logger"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Deps struct {
	RedisClient *redis.Client
}

type Redis struct {
	Deps
	expirationTime time.Duration
}

func New(deps Deps, expiration time.Duration) *Redis {
	return &Redis{
		Deps:           deps,
		expirationTime: expiration,
	}
}

func (c *Redis) Set(ctx context.Context, key string, value string) error {
	return c.RedisClient.Set(ctx, key, value, c.expirationTime).Err()
}

func (c *Redis) Get(ctx context.Context, key string) (string, error) {
	bannerRaw, err := c.RedisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logger.Errorf("redis error: banner id %s does not exist", key)
			return "", err
		}
		logger.Errorf("redis get error: %v", err)
		return "", err
	}
	if bannerRaw == "" {
		logger.Warnf("reddis: banner with id %s is empty", key)
	}
	return bannerRaw, nil
}

func (c *Redis) Delete(ctx context.Context, key string) error {
	err := c.RedisClient.Del(ctx, key).Err()
	if err != nil {
		logger.Errorf("redis delete error: %v", err)
		return err
	}

	logger.Infof("Banner with id %s was deleted from cache", key)
	return nil
}
