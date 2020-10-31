package service

import "math/rand"

type Service interface {
	GetDataByChannel(ch chan<- []byte)
}

func NewService(urls []string, minTimeout int, maxTimeout int, numberOfRequests int) Service {
	return &service{
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
}

func (s *service) GetDataByChannel(ch chan<- []byte) {
	for i := 0; i < s.numberOfRequests; i++ {
		url := s.selectRandomURL()
		ch <- []byte(url)
	}
	close(ch)
}

func (s *service) selectRandomURL() string {
	return s.urls[rand.Intn(len(s.urls))]
}
