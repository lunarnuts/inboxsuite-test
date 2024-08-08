package roadmap_repo

import (
	"github.com/lunarnuts/inboxsuite-test/internal/entity"
	"github.com/lunarnuts/inboxsuite-test/internal/interfaces"
	"github.com/lunarnuts/inboxsuite-test/pkg/postgres"
	"log/slog"
)

var _ interfaces.RoadMapRepository = (*Repository)(nil)

func New(log *slog.Logger) *Repository {
	return &Repository{
		log: log,
	}
}

type Repository struct {
	log *slog.Logger
}

func (r *Repository) LoadCache(db *postgres.Gorm) (map[entity.ClassID]entity.Roadmap, error) {
	const op = "Repository.LoadCache"
	mp := make(map[entity.ClassID]entity.Roadmap)
	var rows []*entity.Row
	if err := db.Model(&entity.Row{}).Find(&rows).Error; err != nil {
		r.log.Error(op, "error", err.Error())
		return nil, err
	}

	for _, row := range rows {
		mp[row.ClassID] = row.Roadmap
	}
	return mp, nil
}
