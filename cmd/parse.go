package cmd

import (
	"fmt"
	"net/url"
	"time"

	"experiments/internal/app"
	crawler "experiments/internal/domain/crawler"
	"experiments/internal/infrastructure/fetcher"
	parser "experiments/internal/infrastructure/parser/html"
	"experiments/internal/usecase/crawling"
	"github.com/spf13/cobra"
)

const (
	depthFlag    = "depth"
	timeoutFlag  = "timeout"
	cooldownFlag = "cooldown"
)

func newParseCommand(app *app.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "parse [url]",
		Short: "omgkotofey go experiments application",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := app.Logger

			parsedURL, err := url.ParseRequestURI(args[0])
			if err != nil {
				err = fmt.Errorf("invalid uri: %w", err)
				logger.Fatal(err.Error())

				return err
			}

			crawlDepth, err := cmd.Flags().GetInt64(depthFlag)
			if err != nil {
				err = fmt.Errorf("invalid depth value: %w", err)
				logger.Fatal(err.Error())

				return err
			}

			timeout, err := cmd.Flags().GetDuration(timeoutFlag)
			if err != nil {
				err = fmt.Errorf("invalid duration value: %w", err)
				logger.Fatal(err.Error())

				return err
			}

			cooldown, err := cmd.Flags().GetDuration(cooldownFlag)
			if err != nil {
				err = fmt.Errorf("invalid cooldown value: %w", err)
				logger.Fatal(err.Error())

				return err
			}

			crawlerInstance := crawling.NewCrawler(
				*app.Config,
				fetcher.NewHTTPFetcher(),
				[]crawler.Parser{
					parser.NewLinksParser(parsedURL),
				},
				logger,
			)

			crawlingResult := crawlerInstance.Crawl(
				crawler.CrawlRequest{
					Context:  cmd.Context(),
					URL:      parsedURL,
					Depth:    crawlDepth,
					Timeout:  timeout,
					Cooldown: cooldown,
				},
			)

			summary := crawlingResult.GetSummary()
			fmt.Println("Results:")
			for i := range summary.GetResults() {
				result := summary.GetResults()[i]
				fmt.Printf(
					"Parsed %s. Response length: %d (%v ms) \n",
					result.GetResource().GetURL(),
					len(result.GetResource().GetBody()),
					result.GetResource().GetResponseTimeMs(),
				)
			}

			if summary.TotalErrors() > 0 {
				fmt.Println("-----------")
				fmt.Println("Errors:")
				for i := range summary.GetErorrs() {
					fmt.Printf("Err: %s\n", summary.GetErorrs()[i])
				}
			}

			fmt.Println("-----------")
			fmt.Printf("Execution Time: %.2f sec\n", summary.GetDuration().Seconds())
			fmt.Printf("Parsed %v urls\n", summary.TotalParsed())
			fmt.Printf("Got %v errors\n", summary.TotalErrors())

			return nil
		},
	}

	var depth int64
	var timeout, cooldown time.Duration
	command.Flags().Int64Var(&depth, depthFlag, -1, "Max depth to crawling (e.g. -1 - unlimited depth, 0 - only specified resource, N - N levels down)")
	command.Flags().DurationVar(&timeout, timeoutFlag, app.Config.Crawler.DefaultFetchTimeout, "Request timeout (e.g. 5s, 1m)")
	command.Flags().DurationVar(&cooldown, cooldownFlag, app.Config.Crawler.DefaultFetchCooldown, "Fetching cooldown (e.g. 300ms, 1s)")

	return command
}
