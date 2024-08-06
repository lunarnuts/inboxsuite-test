package application

import (
	"github.com/lunarnuts/inboxsuite-test/config"
	"github.com/lunarnuts/inboxsuite-test/pkg/logger"
	"github.com/lunarnuts/inboxsuite-test/pkg/postgres"
	"github.com/lunarnuts/inboxsuite-test/pkg/rabbitmq"
	"log/slog"
)

type App struct {
	cfg    *config.Config
	log    *slog.Logger
	db     *postgres.Gorm
	rabbit *rabbitmq.Rabbit
}

func New(cfg *config.Config) (*App, error) {
	const op = "app.New"

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

	return &App{cfg: cfg, log: log, db: db, rabbit: rabbit}, nil
}

func (a *App) Run() error {
	return nil
}
