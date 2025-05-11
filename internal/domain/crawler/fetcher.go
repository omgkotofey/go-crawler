package crawler

import (
	"net/url"
)

type Fetcher interface {
	Fetch(url *url.URL) (result FetchedResource, err error)
}
