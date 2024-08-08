package main

import (
	"github.com/lunarnuts/inboxsuite-test/config"
	"github.com/lunarnuts/inboxsuite-test/internal/entity"
	"github.com/lunarnuts/inboxsuite-test/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func logMigration(db *gorm.DB, table string, action string) {
	log.Printf("Migration %s on table %s", action, table)
}

func main() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DB.ParseURL()), &gorm.Config{
		Logger:         newLogger,
		TranslateError: true,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	if err := db.Callback().Create().Before("gorm:create").Register("before_create", func(tx *gorm.DB) {
		logMigration(tx, "tables", "starting")
	}); err != nil {
		log.Fatalf("Migration failed at before_create callback: %v", err)
	}

	if err := db.Callback().Create().After("gorm:create").Register("after_create", func(tx *gorm.DB) {
		logMigration(tx, "tables", "completed")
	}); err != nil {
		log.Fatalf("Migration failed at after_create callback: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&entity.Row{},
	)
	if err != nil {
		panic(err)
	}

	rows := make([]entity.Row, 1000)
	batchSize := 100

	for i := range rows {
		rows[i] = entity.Row{
			ClassID: entity.ClassID(i),
			Roadmap: entity.Roadmap(i),
		}
	}

	if err = db.CreateInBatches(rows, batchSize).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return
		}
		log.Fatalf("Migration failed at create batches: %v", err)
	}
}
