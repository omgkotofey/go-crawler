package crawl

import (
	"experiments/internal/domain/crawler"
	parsers "experiments/internal/infrastructure/parser/html"
	"fmt"
	"net/url"
	"sync"
)

type Crawler struct {
	wg        *sync.WaitGroup
	fetcher   crawler.Fetcher
	parserSet crawler.ParserSet
	inbox     *crawler.Inbox
}

func NewCrawler(fetcher crawler.Fetcher, parsers []crawler.Parser) Crawler {
	parserSet := crawler.ParserSet{}
	for i := range parsers {
		parserSet.AddParser(parsers[i])
	}

	return Crawler{
		wg:        &sync.WaitGroup{},
		fetcher:   fetcher,
		parserSet: parserSet,
		inbox:     nil,
	}
}

func (c *Crawler) AddParser(parser crawler.Parser) {
	c.parserSet.AddParser(parser)
}

func (c *Crawler) Crawl(urlToCrawl *url.URL, depth int64) (chan string, chan error) {
	resChan := make(chan string)
	errChan := make(chan error)
	c.inbox = crawler.NewInbox()

	c.inbox.Add(urlToCrawl.String(), depth)
	c.wg.Add(1)

	go func() {
		for {
			task, ok := c.inbox.Next()
			if !ok {
				return
			}
			go func() {
				c.crawlUrl(task, resChan, errChan)
				c.wg.Done()
			}()
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

func (c *Crawler) crawlUrl(task crawler.Task, resChan chan string, errChan chan error) {
	if task.Depth < 0 {
		return
	}

	urlToCrawl, err := url.ParseRequestURI(task.URL)
	if err != nil {
		err = fmt.Errorf("invalid url %v: %v", urlToCrawl.String(), err)
		errChan <- err
		return
	}

	fetchResult, err := c.fetcher.Fetch(urlToCrawl)
	if err != nil {
		err = fmt.Errorf("fetching %v: %v", urlToCrawl.String(), err)
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
			c.inbox.Add(urlToCrawl.String(), task.Depth-1)
		}

	}

	resChan <- fmt.Sprintf(
		"fetched %s. response length: %d (%v ms)",
		urlToCrawl,
		len(parseResult.GetResource().GetBody()),
		parseResult.GetResource().GetResponseTimeMs(),
	)
}
