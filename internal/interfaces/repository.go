package interfaces

import (
	"context"
	"github.com/lunarnuts/inboxsuite-test/internal/entity"
	"github.com/lunarnuts/inboxsuite-test/pkg/postgres"
)

type Repository interface {
	Conn() *postgres.Gorm
	ConnWithContext(ctx context.Context) *postgres.Gorm
	RoadMap() RoadMapRepository
}

type RoadMapRepository interface {
	LoadCache(db *postgres.Gorm) (map[entity.ClassID]entity.Roadmap, error)
}
