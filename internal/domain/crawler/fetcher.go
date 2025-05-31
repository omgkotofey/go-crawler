package crawler

import (
	"context"
	"net/url"
)

type Fetcher interface {
	Fetch(ctx context.Context, url *url.URL) (result FetchedResource, err error)
}

type FetchedResource struct {
	url            *url.URL
	body           []byte
	responseTimeMs int64
}

func NewFetchedResource(url *url.URL, body []byte, responseTimeMs int64) FetchedResource {
	return FetchedResource{
		url:            url,
		body:           body,
		responseTimeMs: responseTimeMs,
	}
}

func (r FetchedResource) GetUrl() *url.URL {
	return r.url
}

func (r FetchedResource) GetBody() []byte {
	return r.body
}

func (r FetchedResource) GetResponseTimeMs() int64 {
	return r.responseTimeMs
}
