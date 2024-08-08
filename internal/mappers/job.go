package mappers

import (
	"github.com/lunarnuts/inboxsuite-test/internal/dto"
	"github.com/lunarnuts/inboxsuite-test/internal/entity"
)

func newJobMapper() *ClassIDMapper {
	return &ClassIDMapper{}
}

type ClassIDMapper struct {
}

func (c ClassIDMapper) ToEntity(id dto.ClassID) entity.ClassID {
	return entity.ClassID(id)
}
