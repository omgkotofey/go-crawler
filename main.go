package main

import (
	"errors"
	"experiments/app"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("Execution Time: %.2f sec\n", time.Since(start).Seconds())
	}()

	// args := os.Args
	args := []string{"", "https://go.dev/", "1"}
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
	crawler := app.Crawler{
		Fetcher: app.HttpFetcher{},
		UrlParser: app.TokenizerParser{
			Origin: app.Origin{
				Base: parsedUrl,
			},
		},
	}

	go crawler.Crawl(parsedUrl, crawlDepth, resChan, errChan)

	for {
		select {
		case result, ok := <-resChan:
			fmt.Println(result)
			if !ok {
				resChan = nil
			}
		case err := <-errChan:
			fmt.Println(err)
		}

		if resChan == nil {
			break
		}
	}
}
