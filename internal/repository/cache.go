package repository

import (
	"context"
	"errors"
	"time"

	"github.com/d4niells/shorten/internal/entity"
	"github.com/go-redis/redis"
)

var ErrKeyDoesNotExist = errors.New("key does not exist")

type CacheRepository interface {
	Get(ctx context.Context, key string) (*entity.URL, error)
	Set(ctx context.Context, url *entity.URL, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}

type CacheRepositoryImpl struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *CacheRepositoryImpl {
	return &CacheRepositoryImpl{client}
}

func (c *CacheRepositoryImpl) Get(ctx context.Context, key string) (*entity.URL, error) {
	str, err := c.client.Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrKeyDoesNotExist
		}
		return nil, err
	}

	var url *entity.URL
	err = url.FromJSON(str)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (c *CacheRepositoryImpl) Set(ctx context.Context, url *entity.URL, expiration time.Duration) error {
	str, err := url.ToJSON()
	if err != nil {
		return err
	}

	return c.client.Set(url.Key, str, expiration).Err()
}

func (c *CacheRepositoryImpl) Del(ctx context.Context, key string) error {
	err := c.client.Del(key).Err()
	if err != nil {
		return err
	}

	return nil
}
