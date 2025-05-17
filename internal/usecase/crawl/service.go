package crawl

import (
	"experiments/internal/domain/crawler"
	parsers "experiments/internal/infrastructure/parser/html"
	"fmt"
	"net/url"
	"sync"
)

type urlMap struct {
	urls  map[string]string
	mutex *sync.RWMutex
}

func (m *urlMap) exists(url string) bool {
	m.mutex.Lock()
	_, found := m.urls[url]
	m.mutex.Unlock()

	return found
}

func (m *urlMap) add(url string) {
	m.mutex.Lock()
	m.urls[url] = url
	m.mutex.Unlock()
}

type Crawler struct {
	wg            *sync.WaitGroup
	fetcher       crawler.Fetcher
	parserSet     crawler.ParserSet
	processedUrls *urlMap
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
		processedUrls: &urlMap{
			urls:  make(map[string]string),
			mutex: &sync.RWMutex{},
		},
	}
}

func (c *Crawler) AddParser(parser crawler.Parser) {
	c.parserSet.AddParser(parser)
}

func (c *Crawler) Crawl(urlToCrawl *url.URL, depth int) (chan string, chan error) {
	resChan := make(chan string)
	errChan := make(chan error)

	c.scheduleUrl(urlToCrawl, depth, resChan, errChan)
	go func() {
		c.wg.Wait()
		close(resChan)
		close(errChan)
	}()

	return resChan, errChan
}

func (c *Crawler) scheduleUrl(urlToCrawl *url.URL, depth int, resChan chan string, errChan chan error) {
	c.wg.Add(1)
	go c.crawlUrl(urlToCrawl, depth, resChan, errChan)
}

func (c *Crawler) crawlUrl(urlToCrawl *url.URL, depth int, resChan chan string, errChan chan error) {
	defer c.wg.Done()

	if depth < 0 {
		return
	}

	if c.processedUrls.exists(urlToCrawl.String()) {
		return
	} else {
		c.processedUrls.add(urlToCrawl.String())
	}

	fetchResult, err := c.fetcher.Fetch(urlToCrawl)
	if err != nil {
		err = fmt.Errorf("error occurred while fetching %v: %v", urlToCrawl.String(), err)
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
			c.scheduleUrl(urlToCrawl, depth-1, resChan, errChan)
		}

	}

	resChan <- fmt.Sprintf(
		"fetched %s. response length: %d (%v ms)",
		urlToCrawl,
		len(parseResult.GetResource().GetBody()),
		parseResult.GetResource().GetResponseTimeMs(),
	)
}
