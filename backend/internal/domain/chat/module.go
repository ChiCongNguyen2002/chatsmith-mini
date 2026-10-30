package chat

import (
	"chatsmith-mini/configs"
	"chatsmith-mini/constants"
	"chatsmith-mini/internal/domain/catalog"
	"chatsmith-mini/internal/domain/conversation"
	"chatsmith-mini/pkg/databases"
	"chatsmith-mini/pkg/llm"
)

type Module struct {
	Service *Service
}

func NewModule(db databases.PostgresDB, conversations *conversation.Service, models *catalog.Service, streamers map[constants.Provider]Streamer, classifier *llm.Classifier, cfg configs.ChatConfig) (*Module, error) {
	service, err := newService(newRepo(db), conversations, models, streamers, classifier, cfg)
	if err != nil {
		return nil, err
	}

	return &Module{Service: service}, nil
}
