package postgres

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Gorm struct {
	*gorm.DB
}

func New(connString string, opts ...Option) (*Gorm, error) {
	defaultNowFunc := time.Now

	cfg := &config{
		translateError: true,
		nowFunc:        defaultNowFunc,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	db, err := gorm.Open(postgres.Open(connString), cfg.toGormConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to open connection pool: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.maxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.maxOpenConns)
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping connection: %w", err)
	}
	return &Gorm{db}, nil
}

func (g *Gorm) WithCtx(ctx context.Context) *Gorm {
	return &Gorm{g.WithContext(ctx)}
}

func (g *Gorm) TxBegin(ctx context.Context) *Gorm {
	return &Gorm{g.WithCtx(ctx).Begin()}
}
