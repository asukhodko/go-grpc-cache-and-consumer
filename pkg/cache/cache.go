package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, data []byte) error
}

func NewCache(minTimeoutSec, maxTimeoutSec int, redisAddress string) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})
	return &cache{
		rdb:           rdb,
		minTimeoutSec: minTimeoutSec,
		maxTimeoutSec: maxTimeoutSec,
	}
}

type cache struct {
	rdb           *redis.Client
	minTimeoutSec int
	maxTimeoutSec int
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		data, err = nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "c.rdb.Get")
	}
	return data, nil
}

func (c *cache) Set(ctx context.Context, key string, data []byte) error {
	expSeconds := c.minTimeoutSec + rand.Intn(c.maxTimeoutSec-c.minTimeoutSec)
	exp := time.Duration(expSeconds) * time.Second
	err := c.rdb.Set(ctx, key, data, exp).Err()
	if err != nil {
		return errors.Wrap(err, "c.rdb.Set")
	}
	return nil
}
