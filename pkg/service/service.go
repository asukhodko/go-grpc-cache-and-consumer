package service

import (
	"context"
	"math/rand"
	"sync"

	"github.com/pkg/errors"

	"github.com/asukhodko/go-grps-cache-and-consumer/pkg/cachingfetcher"
)

// Service - интерфейс сервиса
type Service interface {
	GetDataWithChannel(ctx context.Context) (<-chan []byte, <-chan error)
}

// NewService конструирует сервис
func NewService(fetcher cachingfetcher.Fetcher, urls []string, minTimeout int, maxTimeout int, numberOfRequests int) Service {
	return &service{
		fetcher:          fetcher,
		urls:             urls,
		minTimeout:       minTimeout,
		maxTimeout:       maxTimeout,
		numberOfRequests: numberOfRequests,
	}
}

type service struct {
	urls             []string
	minTimeout       int
	maxTimeout       int
	numberOfRequests int
	fetcher          cachingfetcher.Fetcher
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
				data, err := s.fetcher.Fetch(ctx, url)
				if err != nil {
					chErr <- errors.Wrap(err, "s.cachingFetcher.Fetch")
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

func (s *service) selectRandomURL() string {
	return s.urls[rand.Intn(len(s.urls))]
}
