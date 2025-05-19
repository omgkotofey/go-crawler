package cmd

import (
	"errors"
	"experiments/internal/domain/crawler"
	"experiments/internal/infrastructure/fetcher"
	parser "experiments/internal/infrastructure/parser/html"
	"experiments/internal/usecase/crawl"
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
			fmt.Println("-----------")
			fmt.Printf("Execution Time: %.2f sec\n", time.Since(start).Seconds())
			fmt.Printf("Fetched %v urls\n", fetchedUrls)
			fmt.Printf("Got %v errors\n", errorsCount)
		}()

		parsedUrl, err := url.ParseRequestURI(args[0])
		if err != nil {
			panic(err)
		}

		crawlDepth, err := strconv.Atoi(args[1])
		if err != nil || crawlDepth == 0 {
			panic(errors.New("error: invalid depth value"))
		}

		crawler := crawl.NewCrawler(
			fetcher.NewHttpFetcher(),
			[]crawler.Parser{
				parser.NewLinksParser(parsedUrl),
			},
		)

		resChan, errChan := crawler.Crawl(parsedUrl, int64(crawlDepth))

		resChanClosed := false
		errChanClosed := false

		for !resChanClosed || !errChanClosed {
			select {
			case result, ok := <-resChan:
				if !ok {
					resChanClosed = true
					continue
				}
				fmt.Println(result)
				fetchedUrls++
			case err, ok := <-errChan:
				if !ok {
					errChanClosed = true
					continue
				}
				fmt.Println("Err:", err)
				errorsCount++
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(parceCmd)
}
