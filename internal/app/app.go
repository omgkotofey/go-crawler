package app

import (
	"experiments/internal/config"
	"experiments/internal/infrastructure/logger"

	"go.uber.org/zap"
)

type App struct {
	Config *config.Config
	Logger *zap.Logger
}

func New(cfg *config.Config) *App {
	return &App{
		Config: cfg,
		Logger: logger.InitLogger(cfg.IsProduction()),
	}
}
