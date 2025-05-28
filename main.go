package main

import (
	"context"
	"experiments/cmd"
	"experiments/internal/app"
	"experiments/internal/config"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load(ctx)
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
