package service

type Service interface {
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
