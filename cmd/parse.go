package cmd

import (
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

const timeoutFlag = "timeout"
const cooldownFlag = "cooldown"

func newParseCommand(app *app.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "parse [url] [depth]",
		Short: "omgkotofey go experiments application",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
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
				logger.Fatal(err.Error())

				return err
			}

			crawlDepth, err := strconv.Atoi(args[1])
			if err != nil || crawlDepth <= 0 {
				err := errors.New("invalid depth value")
				logger.Fatal(err.Error())

				return err
			}

			timeout, err := cmd.Flags().GetDuration(timeoutFlag)
			if err != nil {
				err := fmt.Errorf("invalid duration value: %w", err)
				logger.Fatal(err.Error())
				return err
			}

			cooldown, err := cmd.Flags().GetDuration(cooldownFlag)
			if err != nil {
				err := fmt.Errorf("invalid cooldown value: %w", err)
				logger.Fatal(err.Error())
				return err
			}

			crawler := crawl.NewCrawler(
				*app.Config,
				fetcher.NewHttpFetcher(),
				[]crawler.Parser{
					parser.NewLinksParser(parsedUrl),
				},
				logger,
			)

			resChan, errChan := crawler.Crawl(
				cmd.Context(),
				crawl.CrawlRequest{
					Url:      parsedUrl,
					Depth:    int64(crawlDepth),
					Timeout:  timeout,
					Cooldown: cooldown,
				},
			)

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

			return nil
		},
	}

	var timeout, cooldown time.Duration
	command.Flags().DurationVar(&timeout, timeoutFlag, app.Config.Crawler.DefaultFetchTimeout, "Request timeout (e.g. 5s, 1m)")
	command.Flags().DurationVar(&cooldown, cooldownFlag, app.Config.Crawler.DefaultFetchCooldown, "Fetching cooldown (e.g. 300ms, 1s)")

	return command
}
