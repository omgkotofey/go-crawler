package main

import (
	"errors"
	"experiments/app/crawler"
	"experiments/app/fetcher"
	"experiments/app/parser"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

// todo: Make tests
// todo: Gorutines pool

func main() {
	start := time.Now()
	fetchedUrls := 0
	errorsCount := 0
	defer func() {
		fmt.Printf("Execution Time: %.2f sec\n", time.Since(start).Seconds())
		fmt.Printf("Fetched %v urls\n", fetchedUrls)
		if errorsCount > 0 {
			fmt.Printf("Got %v errors\n", errorsCount)

		}
	}()

	args := os.Args
	if len(args[1:]) < 2 {
		panic(errors.New("error: invalid input"))
	}

	parsedUrl, err := url.ParseRequestURI(args[1])
	if err != nil {
		panic(err)
	}

	crawlDepth, err := strconv.Atoi(args[2])
	if err != nil || crawlDepth == 0 {
		panic(errors.New("error: invalid depth value"))
	}
	resChan := make(chan string)
	errChan := make(chan error)
	crawler := crawler.Crawler{
		Fetcher: fetcher.NewHttpFetcher(),
		UrlParser: parser.TokenizerParser{
			Origin: parser.Origin{
				Base: parsedUrl,
			},
		},
	}

	go crawler.Crawl(parsedUrl, crawlDepth, resChan, errChan)

	for {
		select {
		case result, ok := <-resChan:
			fmt.Println(result)
			fetchedUrls += 1
			if !ok {
				resChan = nil
			}
		case err := <-errChan:
			fmt.Println(err)
			errorsCount += 1
		}

		if resChan == nil {
			break
		}
	}
}
