package fetcher

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

func (r FetchResult) GetUrl() *url.URL {
	return r.resource
}

func (r FetchResult) GetBody() []byte {
	return r.body
}

func (r FetchResult) GetSpent() float64 {
	return r.spent
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
		resource: u,
		body:     body,
		spent:    time.Since(start).Seconds(),
	}, nil
}

func NewHttpFetcher() HttpFetcher {
	return HttpFetcher{
		Client: &http.Client{},
	}
}
