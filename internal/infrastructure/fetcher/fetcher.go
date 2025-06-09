package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"experiments/internal/domain/crawler"
)

type HTTPFetcher struct {
	Client *http.Client
}

func (f HTTPFetcher) Fetch(ctx context.Context, u *url.URL) (crawler.FetchedResource, error) {
	start := time.Now()
	result := crawler.FetchedResource{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return result, fmt.Errorf("request create: %w", err)
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return result, fmt.Errorf("request send: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("body read: %w", err)
	}

	return crawler.NewFetchedResource(u, body, time.Since(start).Milliseconds()), nil
}

func NewHTTPFetcher() HTTPFetcher {
	return HTTPFetcher{
		Client: &http.Client{},
	}
}
