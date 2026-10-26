package catalog

import (
	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/databases"
)

type CreateModuleDTO struct {
	DB           databases.PostgresDB
	Available    func(constants.Provider) bool
	DefaultModel string
}

type Module struct {
	Service *Service
}

func NewModule(dto *CreateModuleDTO) *Module {
	return &Module{Service: newService(newRepo(dto.DB), dto.Available, dto.DefaultModel)}
}
