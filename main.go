package main

import (
	"context"
	"experiments/cmd"
	"experiments/internal/app"
	"experiments/internal/config"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	command := cmd.NewRootCommand(app.New(cfg))
	if err = command.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	os.Exit(0)
}
