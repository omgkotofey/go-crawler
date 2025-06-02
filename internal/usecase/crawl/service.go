package crawl

import (
	"context"
	"experiments/internal/config"
	"experiments/internal/domain/crawler"
	parsers "experiments/internal/infrastructure/parser/html"
	"fmt"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"
)

type CrawlRequest struct {
	Url      *url.URL
	Depth    int64
	Timeout  time.Duration
	Cooldown time.Duration
}

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

func (c *Crawler) Crawl(ctx context.Context, target CrawlRequest) (chan string, chan error) {
	resChan := make(chan string)
	errChan := make(chan error)
	c.inbox = crawler.NewInbox()
	c.wg = &sync.WaitGroup{}
	c.wg.Add(1)

	defer func() {
		c.inbox.Add(target.Url.String(), target.Depth, target.Timeout)
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
					c.crawlUrl(ctx, task, resChan, errChan)
				}()

				if target.Cooldown != 0 {
					c.logger.Debug(fmt.Sprintf("Cooldown %s", target.Cooldown))
					time.Sleep(target.Cooldown)
				}
			}
		}
	}()

	go func() {
		c.wg.Wait()
		c.inbox.Close()
		close(resChan)
		close(errChan)
	}()

	return resChan, errChan
}

func (c *Crawler) crawlUrl(ctx context.Context, task crawler.Task, resChan chan string, errChan chan error) {
	if task.Depth < 0 {
		return
	}

	urlToCrawl, err := url.ParseRequestURI(task.URL)
	if err != nil {
		err = fmt.Errorf("invalid url %v: %v", urlToCrawl.String(), err)
		errChan <- err
		return
	}

	fetchResult, err := c.fetchUrl(ctx, urlToCrawl, task.Timeout)
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
			urlToCrawl, err := url.ParseRequestURI(parsedLinks[j])
			if err != nil {
				err = fmt.Errorf("invalid url %v: %v", parsedLinks[j], err)
				errChan <- err
				return
			}

			if c.inbox.Exists(urlToCrawl.String()) {
				continue
			}
			c.logger.Debug(fmt.Sprintf("Scheduled %s, depth left: %v", urlToCrawl.String(), task.Depth-1))
			c.wg.Add(1)
			c.inbox.Add(urlToCrawl.String(), task.Depth-1, task.Timeout)
		}

	}

	resChan <- fmt.Sprintf(
		"Fetched %s. response length: %d (%v ms)",
		urlToCrawl,
		len(parseResult.GetResource().GetBody()),
		parseResult.GetResource().GetResponseTimeMs(),
	)
}

func (c *Crawler) fetchUrl(ctx context.Context, urlToFetch *url.URL, timeout time.Duration) (crawler.FetchedResource, error) {
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
		err = fmt.Errorf("fetching %v: %v", urlAsString, err)

		return result, err
	}

	c.logger.Debug(fmt.Sprintf("Finished fetching %s", urlAsString))

	return result, nil
}
