package postgres

import (
	"gorm.io/gorm/logger"
	"time"
)

type config struct {
	translateError bool
	maxIdleConns   int
	maxOpenConns   int
	logMode        logger.LogLevel
	nowFunc        func() time.Time
}
