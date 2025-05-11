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
	fetcher       crawler.Fetcher
	parser        crawler.Parser
	processedUrls *urlMap
}

func NewCrawler(fetcher crawler.Fetcher, parser crawler.Parser) Crawler {
	var mutex *sync.RWMutex
	urlMap := urlMap{
		urls:  make(map[string]string),
		mutex: mutex,
	}

	return Crawler{
		fetcher:       fetcher,
		parser:        parser,
		processedUrls: &urlMap,
	}
}

func (c *Crawler) Crawl(urlToCrawl *url.URL, depth int, resChan chan string, errChan chan error) {
	c.processedUrls = &urlMap{urls: map[string]string{}}

	var wg sync.WaitGroup
	wg.Add(1)
	go c.crawlUrl(urlToCrawl, depth, resChan, errChan, &wg)
	wg.Wait()
	close(resChan)
}

func (c *Crawler) crawlUrl(urlToCrawl *url.URL, depth int, resChan chan string, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()

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
		wg.Add(1)
		go c.crawlUrl(urlToCrawl, depth-1, resChan, errChan, wg)
	}

	resChan <- fmt.Sprintf(
		"fetched %s. response length: %d (%v ms)",
		urlToCrawl,
		len(parseResult.GetResource().GetBody()),
		parseResult.GetResource().GetResponseTimeMs(),
	)
}
