package crawler

import (
	"experiments/internal/domain/crawler"
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
	parser        crawler.Parser
	processedUrls *urlMap
}

func NewCrawler(fetcher crawler.Fetcher, parser crawler.Parser) Crawler {
	return Crawler{
		wg:      &sync.WaitGroup{},
		fetcher: fetcher,
		parser:  parser,
		processedUrls: &urlMap{
			urls:  make(map[string]string),
			mutex: &sync.RWMutex{},
		},
	}
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

	parseResult, err := c.parser.Parse(fetchResult)
	if err != nil {
		err = fmt.Errorf("error occurred while parsing %v: %v", urlToCrawl.String(), err)
		errChan <- err
		return
	}

	for i := range parseResult.GetData() {
		urlToCrawl, err := url.ParseRequestURI(parseResult.GetData()[i])
		if err != nil {
			err = fmt.Errorf("invalid url %v: %v", parseResult.GetData()[i], err)
			errChan <- err
			return
		}
		c.scheduleUrl(urlToCrawl, depth-1, resChan, errChan)
	}

	resChan <- fmt.Sprintf(
		"fetched %s. response length: %d (%v ms)",
		urlToCrawl,
		len(parseResult.GetResource().GetBody()),
		parseResult.GetResource().GetResponseTimeMs(),
	)
}
