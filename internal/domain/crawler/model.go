package crawler

import (
	"net/url"
)

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

type ParsedData struct {
	resource FetchedResource
	data     []string
}

func NewParsedData(resource FetchedResource, data []string) ParsedData {
	return ParsedData{resource: resource, data: data}
}

func (p ParsedData) GetResource() FetchedResource {
	return p.resource
}

func (p ParsedData) GetData() []string {
	return p.data
}

func (p ParsedData) AppendData(rows []string) {
	p.data = append(p.data, rows...)
}
