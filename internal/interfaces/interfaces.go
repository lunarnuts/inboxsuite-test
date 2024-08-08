package interfaces

import (
	"context"
	"github.com/lunarnuts/inboxsuite-test/internal/entity"
)

type Runner interface {
	Run(ctx context.Context)
}

type Notifier interface {
	Notify() <-chan error
}

type RunnerNotifier interface {
	Runner
	Notifier
}

type CacheLoader interface {
	RoadmapGetter
}

type RoadmapGetter interface {
	Get(id entity.ClassID) (entity.Roadmap, error)
}
