package fetcher

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

type FetchResult struct {
	Resource *url.URL
	Body     []byte
	Spent    float64
}

type Fetcher interface {
	Fetch(url *url.URL) (result FetchResult, err error)
}

type HttpFetcher struct {
	Client *http.Client
}

func (f HttpFetcher) Fetch(u *url.URL) (result FetchResult, err error) {
	start := time.Now()
	resp, err := f.Client.Get(u.String())
	if err != nil {
		return FetchResult{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FetchResult{}, err
	}

	return FetchResult{
		Resource: u,
		Body:     body,
		Spent:    time.Since(start).Seconds(),
	}, nil
}

func NewHttpFetcher() HttpFetcher {
	return HttpFetcher{
		Client: &http.Client{},
	}
}
