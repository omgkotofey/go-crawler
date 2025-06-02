package config

import (
	"context"
	"fmt"
	"time"

	"github.com/sethvargo/go-envconfig"
)

const production = "production"
const development = "development"

type AppConfig struct {
	Env     string `env:"ENV,default=development"`
	Version string `env:"VERSION,required"`
}

type CrawlerConfig struct {
	MaxParallelFetches   int           `env:"MAX_PARALLEL_FETCHES,default=100"`
	DefaultFetchTimeout  time.Duration `env:"DEFAULT_FETCH_TIMEOUT,default=3s"`
	DefaultFetchCooldown time.Duration `env:"DEFAULT_FETCH_COOLDOWN,default=0ms"`
}

type Config struct {
	App     AppConfig     `env:", prefix=APP_"`
	Crawler CrawlerConfig `env:", prefix=CRAWLER_"`
}

func (cfg *Config) Validate() {
	if cfg.App.Env != production && cfg.App.Env != development {
		panic(fmt.Sprintf("Unknown environment \"%s\"", cfg.App.Env))
	}
}

func (cfg *Config) IsProduction() bool {
	return cfg.App.Env == production
}

func Load(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}
	cfg.Validate()

	return &cfg, nil
}
