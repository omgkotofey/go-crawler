package cmd

import (
	"errors"
	"experiments/internal/infrastructure/fetcher"
	parser "experiments/internal/infrastructure/parser/html"
	"experiments/internal/usecase/crawler"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var parceCmd = &cobra.Command{
	Use:   "parse [url] [depth]",
	Short: "omgkotofey go experiments application",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
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

		parsedUrl, err := url.ParseRequestURI(args[0])
		if err != nil {
			panic(err)
		}

		crawlDepth, err := strconv.Atoi(args[1])
		if err != nil || crawlDepth == 0 {
			panic(errors.New("error: invalid depth value"))
		}
		resChan := make(chan string)
		errChan := make(chan error)
		crawler := crawler.NewCrawler(fetcher.NewHttpFetcher(), parser.NewHTMLParser(parsedUrl))

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
	},
}

func init() {
	rootCmd.AddCommand(parceCmd)
}
