package postgres

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

type Option func(*config)

func NowFunc(f func() time.Time) Option {
	return func(c *config) {
		c.nowFunc = f
	}
}

func MaxIdleConns(n int) Option {
	return func(c *config) {
		c.maxIdleConns = n
	}
}

func MaxOpenConns(n int) Option {
	return func(c *config) {
		c.maxOpenConns = n
	}
}

func SilentLogger() Option {
	return func(c *config) {
		c.logMode = logger.Silent
	}
}

func LogLevel(level string) Option {
	level = strings.ToLower(level)
	intLevel := logger.Info
	switch level {
	case "info":
		intLevel = logger.Info
	case "warn":
		intLevel = logger.Warn
	case "error":
		intLevel = logger.Error
	default:
		intLevel = logger.Silent
	}
	return func(c *config) {
		c.logMode = intLevel
	}
}

func (c *config) toGormConfig() *gorm.Config {
	return &gorm.Config{
		NowFunc:        c.nowFunc,
		TranslateError: c.translateError,
		Logger:         logger.Default.LogMode(c.logMode),
	}
}
