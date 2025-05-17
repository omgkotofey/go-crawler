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

type ParsedResource struct {
	resource      FetchedResource
	parcingResult []ParsedData
}

func (r ParsedResource) GetResource() FetchedResource {
	return r.resource
}

func (p *ParsedResource) AddResults(results []ParsedData) {
	p.parcingResult = append(p.parcingResult, results...)
}

func (p *ParsedResource) GetResults() []ParsedData {
	return p.parcingResult
}

func NewParsedResource(resource FetchedResource, results []ParsedData) ParsedResource {
	return ParsedResource{resource: resource, parcingResult: results}
}

type ParsedData struct {
	parser ParserType
	data   []string
	err    error
}

func NewParsedData(parser ParserType, data []string, err error) ParsedData {
	return ParsedData{parser: parser, data: data, err: err}
}

func (p *ParsedData) IsSuccess() bool {
	return p.err == nil
}

func (p *ParsedData) GetParserType() ParserType {
	return p.parser
}

func (p *ParsedData) GetData() []string {
	return p.data
}

func (p *ParsedData) AppendData(rows []string) {
	p.data = append(p.data, rows...)
}

func (p *ParsedData) SetError(err error) {
	p.err = err
}
