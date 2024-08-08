package main

import (
	"context"
	"github.com/lunarnuts/inboxsuite-test/config"
	"github.com/lunarnuts/inboxsuite-test/internal/application"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	appCtx := context.Background()

	app, err := application.New(appCtx, cfg)
	if err != nil {
		os.Exit(1)
	}

	err = app.InitCache()
	if err != nil {
		os.Exit(1)
	}

	err = app.InitWorkers()
	if err != nil {
		os.Exit(1)
	}

	if err = app.Run(); err != nil {
		os.Exit(1)
	}
}
