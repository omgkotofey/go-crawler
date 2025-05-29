package main

import (
	"context"
	"experiments/cmd"
	"experiments/internal/app"
	"experiments/internal/config"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	command := cmd.NewRootCommand(app.New(cfg))
	if err = command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
