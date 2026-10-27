package chat

import (
	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/catalog"
	"chatsmith-mini/internal/conversation"
	"chatsmith-mini/pkg/databases"
)

type CreateModuleDTO struct {
	DB            databases.PostgresDB
	Catalog       *catalog.Service
	Conversations *conversation.Service
	Streamers     map[constants.Provider]Streamer
	Config        configs.ChatConfig
}

type Module struct {
	Service *Service
}

func NewModule(dto *CreateModuleDTO) *Module {
	return &Module{Service: newService(newRepo(dto.DB), dto.Catalog, dto.Conversations, dto.Streamers, dto.Config)}
}
