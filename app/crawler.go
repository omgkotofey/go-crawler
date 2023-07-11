package app

import (
	"fmt"
	"net/url"
	"sync"
)

type urlMap map[string]string

func (m urlMap) exists(url string) bool {
	_, found := m[url]

	return found
}

func (m urlMap) add(url string) {
	m[url] = url
}

type Crawler struct {
	Fetcher
	UrlParser
	fetchedUrls urlMap
}

func (c *Crawler) Crawl(urlToCrawl *url.URL, depth int, resChan chan string, errChan chan error) {

	c.fetchedUrls = urlMap{}

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

	if c.fetchedUrls.exists(urlToCrawl.String()) {
		return
	}

	fetchResult, err := c.Fetcher.Fetch(urlToCrawl)
	c.fetchedUrls.add(urlToCrawl.String())
	if err != nil {
		err = fmt.Errorf("error occurred while fetching %v: %v", urlToCrawl.String(), err)
		errChan <- err
		return
	}

	parseResult, err := c.UrlParser.Parse(fetchResult)
	if err != nil {
		err = fmt.Errorf("error occurred while parsing %v: %v", urlToCrawl.String(), err)
		errChan <- err
		return
	}

	wg.Add(len(parseResult.urls))
	for _, urlToCrawl := range parseResult.urls {
		go c.crawlUrl(urlToCrawl, depth-1, resChan, errChan, wg)
	}

	resChan <- fmt.Sprintf("fetched %s. response length: %d (%.2f sec)", urlToCrawl, len(parseResult.body), parseResult.spent)
}
