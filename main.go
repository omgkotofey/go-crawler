package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"experiments/cmd"
	"experiments/internal/app"
	"experiments/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to .env file: %v", err))
	}

	cfg, err := config.Load(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	command := cmd.NewRootCommand(app.New(cfg))
	if err = command.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
