package catalog

import (
	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/databases"
)

type Module struct {
	Service *Service
}

func NewModule(db databases.PostgresDB, available func(constants.Provider) bool, defaultModel string) *Module {
	return &Module{Service: newService(newRepo(db), available, defaultModel)}
}
