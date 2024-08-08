package cache

import (
	"context"
	"github.com/lunarnuts/inboxsuite-test/internal/entity"
	"github.com/lunarnuts/inboxsuite-test/internal/interfaces"
	"github.com/lunarnuts/inboxsuite-test/pkg/errors"
)

var _ interfaces.CacheLoader = (*Cache)(nil)

func New(ctx context.Context, repository interfaces.Repository) (*Cache, error) {
	cache, err := repository.RoadMap().LoadCache(repository.ConnWithContext(ctx))
	return &Cache{
		cache: cache,
	}, err
}

type Cache struct {
	cache map[entity.ClassID]entity.Roadmap
}

func (c Cache) Get(id entity.ClassID) (entity.Roadmap, error) {
	val, ok := c.cache[id]
	if !ok {
		return 0, errors.New("not found")
	}
	return val, nil
}
