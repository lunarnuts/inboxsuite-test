package roadmap_repo

import "github.com/lunarnuts/inboxsuite-test/internal/interfaces"

var _ interfaces.RoadMapRepository = (*Repository)(nil)

func New() *Repository {
	return &Repository{}
}

type Repository struct {
}
