package user

import "chatsmith-mini/pkg/databases"

type CreateModuleDTO struct {
	DB databases.PostgresDB
}

type Module struct {
	Service *Service
}

func NewModule(dto *CreateModuleDTO) *Module {
	return &Module{Service: newService(newRepo(dto.DB))}
}
