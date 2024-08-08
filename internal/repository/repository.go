package repository

import (
	"context"
	"github.com/lunarnuts/inboxsuite-test/internal/interfaces"
	roadmap_repo "github.com/lunarnuts/inboxsuite-test/internal/repository/roadmap"
	"github.com/lunarnuts/inboxsuite-test/pkg/postgres"
	"log/slog"
)

var _ interfaces.Repository = (*Repository)(nil)

func New(db *postgres.Gorm, log *slog.Logger) *Repository {
	return &Repository{
		roadmap: roadmap_repo.New(log),
		db:      db,
		log:     log,
	}
}

type Repository struct {
	roadmap interfaces.RoadMapRepository
	db      *postgres.Gorm
	log     *slog.Logger
}

func (r *Repository) RoadMap() interfaces.RoadMapRepository {
	return r.roadmap
}

func (r *Repository) Conn() *postgres.Gorm {
	return r.db
}

func (r *Repository) ConnWithContext(ctx context.Context) *postgres.Gorm {
	return r.Conn().WithCtx(ctx)
}
