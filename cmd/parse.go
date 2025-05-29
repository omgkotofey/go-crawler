package cmd

import (
	"context"
	"errors"
	"experiments/internal/app"
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

func newParseCommand(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "parse [url] [depth]",
		Short: "omgkotofey go experiments application",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			start := time.Now()
			fetchedUrls := 0
			errorsCount := 0
			logger := app.Logger

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
				*app.Config,
				fetcher.NewHttpFetcher(),
				[]crawler.Parser{
					parser.NewLinksParser(parsedUrl),
				},
				logger,
			)

			resChan, errChan := crawler.Crawl(context.Background(), parsedUrl, int64(crawlDepth))

			resChanClosed := false
			errChanClosed := false

			for !resChanClosed || !errChanClosed {
				select {
				case result, ok := <-resChan:
					if !ok {
						resChanClosed = true
						continue
					}
					logger.Info(result)
					fetchedUrls++
				case err, ok := <-errChan:
					if !ok {
						errChanClosed = true
						continue
					}
					logger.Error(fmt.Sprintf("Err: %s", err))
					errorsCount++
				}
			}

		},
	}
}
