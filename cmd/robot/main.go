package main

import (
	"context"
	"log"
	"os"

	"github.com/undeadpelmen/webrobot-robot/internal/app"
)

func main() {
	application := app.NewApplication()
	flags := application.ParseFlags()

	if err := application.Initialize(flags); err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	ctx := context.Background()
	if err := application.Run(ctx); err != nil {
		log.Fatal("Application failed:", err)
	}

	os.Exit(0)
}
