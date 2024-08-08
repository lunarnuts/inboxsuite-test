package application

import (
	"context"
	"github.com/lunarnuts/inboxsuite-test/config"
	"github.com/lunarnuts/inboxsuite-test/internal/cache"
	"github.com/lunarnuts/inboxsuite-test/internal/interfaces"
	"github.com/lunarnuts/inboxsuite-test/internal/pool"
	"github.com/lunarnuts/inboxsuite-test/internal/repository"
	"github.com/lunarnuts/inboxsuite-test/pkg/logger"
	"github.com/lunarnuts/inboxsuite-test/pkg/postgres"
	"github.com/lunarnuts/inboxsuite-test/pkg/rabbitmq"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	ctx      context.Context
	cfg      *config.Config
	log      *slog.Logger
	db       *postgres.Gorm
	rabbit   *rabbitmq.Rabbit
	shutdown chan os.Signal
	pool     interfaces.RunnerNotifier
	cache    interfaces.CacheLoader
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	const op = "main.New"

	log := logger.SetupLogger(logger.Config{
		Level: logger.Level(cfg.Logger.Level),
		Env:   cfg.Logger.Env,
	})

	db, err := postgres.New(cfg.DB.ParseURL(),
		postgres.MaxIdleConns(cfg.DB.MaxIdleConnections),
		postgres.MaxOpenConns(cfg.DB.MaxOpenConnections),
		postgres.LogLevel(cfg.DB.LogLevel))
	if err != nil {
		log.Error(op, "error", err.Error())
		return nil, err
	}

	rabbit, err := rabbitmq.New(cfg.RabbitMQ.ParseURL())
	if err != nil {
		log.Error(op, "error", err.Error())
		return nil, err
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	return &App{
		cfg:      cfg,
		log:      log,
		db:       db,
		rabbit:   rabbit,
		shutdown: shutdown,
		ctx:      ctx}, nil
}

func (a *App) InitCache() error {
	const op = "main.InitCache"
	var err error
	a.cache, err = cache.New(a.ctx,
		repository.New(a.db, a.log))
	if err != nil {
		a.log.Error(op, "error", err.Error())
	}
	return nil
}

func (a *App) InitWorkers() error {
	const op = "main.InitWorkers"
	var err error
	a.pool, err = pool.New(a.ctx, a.cfg, a.cache, a.log, a.rabbit)
	if err != nil {
		a.log.Error(op, "error", err.Error())
	}
	return nil
}

func (a *App) Run() error {
	const op = "main.Run"
	_, cancel := context.WithCancel(a.ctx)
	a.log.Info("Starting application...")
	defer a.log.Info("Application stopped.")

	go a.pool.Run(a.ctx)

	select {
	case err := <-a.pool.Notify():
		a.log.Error(op, "worker pool error", err.Error())
	case signal := <-a.shutdown:
		a.log.Info(op, "signal", signal)
	}
	cancel()
	return nil
}

func (a *App) stop() {
	close(a.shutdown)

	a.log.Info("Application stopped")
}
