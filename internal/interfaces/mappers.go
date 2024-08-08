package interfaces

import (
	"github.com/lunarnuts/inboxsuite-test/internal/dto"
	"github.com/lunarnuts/inboxsuite-test/internal/entity"
)

type Mappers interface {
	ClassID() ClassIDMapper
	RoadmapID() RoadmapIDMapper
}

type RoadmapIDMapper interface {
	ToDTO(id entity.Roadmap) dto.RoadmapID
}

type ClassIDMapper interface {
	ToEntity(id dto.ClassID) entity.ClassID
}
