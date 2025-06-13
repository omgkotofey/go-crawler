package crawling

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"experiments/internal/config"
	"experiments/internal/domain/crawler"
	parsers "experiments/internal/infrastructure/parser/html"
	"go.uber.org/zap"
)

type Crawler struct {
	wg        *sync.WaitGroup
	limiter   chan struct{}
	fetcher   crawler.Fetcher
	parserSet crawler.ParserSet
	inbox     *crawler.Inbox
	logger    *zap.Logger
}

func NewCrawler(cfg config.Config, fetcher crawler.Fetcher, parsers []crawler.Parser, logger *zap.Logger) Crawler {
	parserSet := crawler.ParserSet{}
	for i := range parsers {
		parserSet.AddParser(parsers[i])
	}

	limit := cfg.Crawler.MaxParallelFetches
	if limit == 0 {
		limit = 1
	}
	limiter := make(chan struct{}, limit)

	return Crawler{
		wg:        nil,
		limiter:   limiter,
		fetcher:   fetcher,
		parserSet: parserSet,
		logger:    logger,
		inbox:     nil,
	}
}

func (c *Crawler) AddParser(parser crawler.Parser) {
	c.parserSet.AddParser(parser)
}

func (c *Crawler) Crawl(request crawler.CrawlRequest) *crawler.CrawlResult {
	ctx := request.Context
	result := crawler.NewCrawlResultForRequest(request)
	resChan, errChan := result.Channels()

	c.inbox = crawler.NewInbox()
	c.wg = &sync.WaitGroup{}
	c.wg.Add(1)

	defer func() {
		c.inbox.Add(request.URL.String(), request.Depth, request.Timeout)
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.logger.Warn("Context cancelled")
				return
			default:
				task, ok := c.inbox.Next()
				if !ok {
					return
				}

				go func() {
					defer c.wg.Done()
					c.crawlURL(ctx, task, resChan, errChan)
				}()

				if request.Cooldown != 0 {
					c.logger.Debug(fmt.Sprintf("Cooldown %s", request.Cooldown))
					time.Sleep(request.Cooldown)
				}
			}
		}
	}()

	go func() {
		c.wg.Wait()
		c.inbox.Close()
		result.Done()
	}()

	return result
}

func (c *Crawler) crawlURL(ctx context.Context, task crawler.Task, resChan chan crawler.ParsedResource, errChan chan error) {
	urlToCrawl, err := url.Parse(task.URL)
	if err != nil {
		err = fmt.Errorf("invalid url %v: %w", urlToCrawl.String(), err)
		errChan <- err
		return
	}

	fetchResult, err := c.fetchURL(ctx, urlToCrawl, task.Timeout)
	if err != nil {
		errChan <- err

		return
	}

	parseResult := c.parserSet.Parse(fetchResult)
	for i := range parseResult.GetResults() {
		if parseResult.GetResults()[i].GetParserType() != parsers.LinksParserName {
			continue
		}

		parsedLinks := parseResult.GetResults()[i].GetData()
		for j := range parsedLinks {
			urlToCrawl, err := url.Parse(parsedLinks[j])
			if err != nil {
				err = fmt.Errorf("invalid url %v: %w", parsedLinks[j], err)
				errChan <- err
				return
			}

			if c.inbox.Exists(urlToCrawl.String()) {
				c.logger.Debug(fmt.Sprintf("Skipped %s: already seen", urlToCrawl.String()))
				continue
			}

			depthLeft := task.Depth
			if task.Depth > 0 {
				depthLeft = task.Depth - 1
			}

			if depthLeft == 0 {
				return
			}

			c.logger.Debug(fmt.Sprintf("Scheduled %s, depth left: %v", urlToCrawl.String(), depthLeft))
			c.wg.Add(1)
			c.inbox.Add(urlToCrawl.String(), task.Depth-1, task.Timeout)
		}
	}

	resChan <- parseResult
}

func (c *Crawler) fetchURL(ctx context.Context, urlToFetch *url.URL, timeout time.Duration) (crawler.FetchedResource, error) {
	urlAsString := urlToFetch.String()
	var result crawler.FetchedResource

	c.limiter <- struct{}{}
	defer func() {
		<-c.limiter
	}()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	c.logger.Debug(fmt.Sprintf("Start fetching %s", urlAsString))

	result, err := c.fetcher.Fetch(ctxWithTimeout, urlToFetch)
	if err != nil {
		err = fmt.Errorf("fetching %v: %w", urlAsString, err)

		return result, err
	}

	c.logger.Debug(fmt.Sprintf("Finished fetching %s", urlAsString))

	return result, nil
}
