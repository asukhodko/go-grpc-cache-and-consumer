package cache

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Cache interface {
	GetOrSetWhenNotExists(ctx context.Context, key string, f func() ([]byte, error)) ([]byte, error)
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
	mutex         sync.Mutex
	keyMutexes    sync.Map
}

func (c *cache) GetOrSetWhenNotExists(ctx context.Context, key string, f func() ([]byte, error)) ([]byte, error) {
	c.mutex.Lock()
	m, ok := c.keyMutexes.Load(key)
	if !ok {
		m = &sync.Mutex{}
		c.keyMutexes.Store(key, m)
	}
	c.mutex.Unlock()

	m.(*sync.Mutex).Lock()
	defer m.(*sync.Mutex).Unlock()

	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err != redis.Nil {
			log.Printf("[WARN] error get from cache for %s: %v", key, err)
		}
		data, err = f()
		if err != nil {
			return nil, errors.Wrap(err, "f")
		}

		expSeconds := c.minTimeoutSec + rand.Intn(c.maxTimeoutSec-c.minTimeoutSec)
		exp := time.Duration(expSeconds) * time.Second
		err := c.rdb.Set(ctx, key, data, exp).Err()
		if err != nil {
			return nil, errors.Wrap(err, "c.rdb.Set")
		}

		return data, nil
	}

	return data, nil
}
