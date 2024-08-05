package repository

import "github.com/lunarnuts/inboxsuite-test/internal/interfaces"

var _ interfaces.Repository = (*Repository)(nil)

func New() *Repository {
	return &Repository{
		roadmap: roadmap_repo.New(),
	}
}

type Repository struct {
	roadmap interfaces.RoadMapRepository
}

func (r Repository) RoadMap() interfaces.RoadMapRepository {
	return r.roadmap
}
