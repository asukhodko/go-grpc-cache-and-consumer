package service

import (
	"context"
	"log"
	"math/rand"
	"sync"

	"github.com/pkg/errors"

	"github.com/asukhodko/go-grpc-cache-and-consumer/pkg/cache"
	"github.com/asukhodko/go-grpc-cache-and-consumer/pkg/urlfetcher"
)

// Service - интерфейс сервиса
type Service interface {
	GetDataWithinChannel(ctx context.Context) (<-chan []byte, <-chan error)
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

func (s *service) GetDataWithinChannel(ctx context.Context) (<-chan []byte, <-chan error) {
	chData := make(chan []byte)
	chErr := make(chan error)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(s.numberOfRequests)
		for i := 0; i < s.numberOfRequests; i++ {
			go func() {
				defer wg.Done()
				data, err := s.getNextData(ctx)
				if err != nil {
					chErr <- errors.Wrap(err, "s.getNextData")
					return
				}
				chData <- data
			}()
		}
		wg.Wait()
		close(chData)
	}()

	return chData, chErr
}

func (s *service) getNextData(ctx context.Context) ([]byte, error) {
	url := s.selectRandomURL()
	data, err := s.cache.GetOrSetWhenNotExists(ctx, url, func() ([]byte, error) {
		log.Printf("[DEBUG] cache miss: %s", url)
		data, err := s.fetcher.Fetch(ctx, url)
		if err != nil {
			return nil, errors.Wrap(err, "s.fetcher.Fetch")
		}
		return data, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "s.cache.GetOrSetWhenNotExists")
	}
	return data, nil
}

func (s *service) selectRandomURL() string {
	return s.urls[rand.Intn(len(s.urls))]
}
