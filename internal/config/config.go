package config

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

const production = "production"
const development = "development"

type appConfig struct {
	Debug bool   `env:"DEBUG,default=false"`
	Env   string `env:"ENV,default=development"`
}

type crawlerConfig struct {
	WorkerPool int `env:"WORKER_POOL,default=100"`
}

type Config struct {
	App     appConfig     `env:"APP_"`
	Crawler crawlerConfig `env:"CRAWLER_"`
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
