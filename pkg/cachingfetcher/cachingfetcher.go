package cachingfetcher

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

// Fetcher - интерфейс получения документа по URL
type Fetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
}

// NewFetcher конструирует фетчер
func NewFetcher() Fetcher {
	return &fetcher{}
}

type fetcher struct{}

func (f *fetcher) Fetch(_ context.Context, url string) ([]byte, error) {
	status, body, err := fasthttp.Get(nil, url)
	if err != nil {
		return nil, errors.Wrapf(err, "fasthttp.Do (url=%s)", url)
	}
	if status >= http.StatusBadRequest {
		return nil, errors.Errorf("unsuccessful status code [%d] from url %s", status, url)
	}

	return body, nil
}
