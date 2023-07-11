package app

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

type FetchResult struct {
	resource *url.URL
	body     []byte
	spent    float64
}

type Fetcher interface {
	Fetch(url *url.URL) (result FetchResult, err error)
}

type HttpFetcher struct{}

func (s HttpFetcher) Fetch(u *url.URL) (result FetchResult, err error) {
	start := time.Now()
	resp, err := http.Get(u.String())
	if err != nil {
		return FetchResult{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FetchResult{}, err
	}

	return FetchResult{
		resource: u,
		body:     body,
		spent:    time.Since(start).Seconds(),
	}, nil
}
