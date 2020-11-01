package service

import (
	"context"
	"log"
	"math/rand"
	"sync"

	"github.com/pkg/errors"

	"github.com/asukhodko/go-grps-cache-and-consumer/pkg/cache"
	"github.com/asukhodko/go-grps-cache-and-consumer/pkg/urlfetcher"
)

// Service - интерфейс сервиса
type Service interface {
	GetDataWithChannel(ctx context.Context) (<-chan []byte, <-chan error)
}

// NewService конструирует сервис
func NewService(fetcher urlfetcher.Fetcher, cache cache.Cache, urls []string, numberOfRequests int) Service {
	return &service{
		fetcher:          fetcher,
		cache:            cache,
		urls:             urls,
		numberOfRequests: numberOfRequests,
	}
}

type service struct {
	fetcher          urlfetcher.Fetcher
	cache            cache.Cache
	urls             []string
	numberOfRequests int
}

func (s *service) GetDataWithChannel(ctx context.Context) (<-chan []byte, <-chan error) {
	chData := make(chan []byte)
	chErr := make(chan error)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(s.numberOfRequests)
		for i := 0; i < s.numberOfRequests; i++ {
			go func() {
				defer wg.Done()
				url := s.selectRandomURL()
				data, err := s.cache.Get(ctx, url)
				if err != nil {
					log.Printf("[WARN] error get from cache for %s: %v", url, err)
					data, err = nil, nil
				}
				if len(data) == 0 {
					data, err = s.fetcher.Fetch(ctx, url)
					if err != nil {
						chErr <- errors.Wrap(err, "s.cachingFetcher.Fetch")
						return
					}
					err = s.cache.Set(ctx, url, data)
					if err != nil {
						log.Printf("[WARN] error set to cache for %s: %v", url, err)
						err = nil
					}
				}
				chData <- data
			}()
		}
		wg.Wait()
		close(chData)
	}()

	return chData, chErr
}

func (s *service) selectRandomURL() string {
	return s.urls[rand.Intn(len(s.urls))]
}
