package mappers

import "github.com/lunarnuts/inboxsuite-test/internal/interfaces"

var _ interfaces.Mappers = (*Mapper)(nil)

func New() *Mapper {
	return &Mapper{
		jobMapper:     newJobMapper(),
		roadmapMapper: newRoadmapMapper(),
	}
}

type Mapper struct {
	jobMapper     interfaces.ClassIDMapper
	roadmapMapper interfaces.RoadmapIDMapper
}

func (m Mapper) RoadmapID() interfaces.RoadmapIDMapper {
	return m.roadmapMapper
}

func (m Mapper) ClassID() interfaces.ClassIDMapper {
	return m.jobMapper
}
