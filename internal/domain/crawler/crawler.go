package crawler

import (
	"context"
	"net/url"
	"time"
)

type CrawlRequest struct {
	Context  context.Context
	URL      *url.URL
	Depth    int64
	Timeout  time.Duration
	Cooldown time.Duration
	Limit    int64
}

type CrawlingSummary struct {
	results  []ParsedResource
	errors   []error
	duration time.Duration
}

func (c CrawlingSummary) TotalParsed() int64 {
	return int64(len(c.results))
}

func (c CrawlingSummary) TotalErrors() int64 {
	return int64(len(c.errors))
}

func (c CrawlingSummary) GetResults() []ParsedResource {
	return c.results
}

func (c CrawlingSummary) GetErorrs() []error {
	return c.errors
}

func (c CrawlingSummary) GetDuration() time.Duration {
	return c.duration
}

type CrawlResult struct {
	done    bool
	request CrawlRequest
	results chan ParsedResource
	errors  chan error
	start   time.Time
	end     time.Time
	summary *CrawlingSummary
}

func NewCrawlResultForRequest(request CrawlRequest) *CrawlResult {
	resChan := make(chan ParsedResource)
	errChan := make(chan error)

	return &CrawlResult{
		request: request,
		results: resChan,
		errors:  errChan,
		start:   time.Now(),
		summary: nil,
	}
}

func (c *CrawlResult) Channels() (chan ParsedResource, chan error) {
	return c.results, c.errors
}

func (c *CrawlResult) Request() CrawlRequest {
	return c.request
}

func (c *CrawlResult) IsDone() bool {
	return c.done
}

func (c *CrawlResult) IsInProgress() bool {
	return !c.done
}

func (c *CrawlResult) GetSummary() *CrawlingSummary {
	if c.summary != nil {
		return c.summary
	}

	results := make([]ParsedResource, 0)
	errors := make([]error, 0)

	for c.IsInProgress() {
		select {
		case result, ok := <-c.results:
			if !ok {
				continue
			}

			results = append(results, result)
		case err, ok := <-c.errors:
			if !ok {
				continue
			}

			errors = append(errors, err)
		}
	}

	c.summary = &CrawlingSummary{
		results:  results,
		errors:   errors,
		duration: c.end.Sub(c.start),
	}

	return c.summary
}

func (c *CrawlResult) Done() {
	if c.done {
		return
	}

	close(c.errors)
	close(c.results)
	c.done = true
	c.end = time.Now()
}

type Crawler interface {
	Crawl(target CrawlRequest) CrawlResult
}
