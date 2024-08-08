package mappers

import (
	"github.com/lunarnuts/inboxsuite-test/internal/dto"
	"github.com/lunarnuts/inboxsuite-test/internal/entity"
	"github.com/lunarnuts/inboxsuite-test/internal/interfaces"
)

func newRoadmapMapper() interfaces.RoadmapIDMapper {
	return &RoadmapMapper{}
}

type RoadmapMapper struct {
}

func (r RoadmapMapper) ToDTO(id entity.Roadmap) dto.RoadmapID {
	return dto.RoadmapID(id)
}
