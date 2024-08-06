package main

import (
	"github.com/lunarnuts/inboxsuite-test/config"
	"github.com/lunarnuts/inboxsuite-test/internal/application"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	app, err := application.New(cfg)
	if err != nil {
		os.Exit(1)
	}

	if err = app.Run(); err != nil {
		os.Exit(1)
	}
}
