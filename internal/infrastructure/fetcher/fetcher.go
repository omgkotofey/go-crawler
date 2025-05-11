package fetcher

import (
	"experiments/internal/domain/crawler"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HttpFetcher struct {
	Client *http.Client
}

func (f HttpFetcher) Fetch(u *url.URL) (crawler.FetchedResource, error) {
	start := time.Now()
	result := crawler.FetchedResource{}
	resp, err := f.Client.Get(u.String())
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	return crawler.NewFetchedResource(u, body, time.Since(start).Milliseconds()), nil
}

func NewHttpFetcher() HttpFetcher {
	return HttpFetcher{
		Client: &http.Client{},
	}
}
