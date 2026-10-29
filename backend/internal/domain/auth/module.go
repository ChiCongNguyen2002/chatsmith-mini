package auth

import (
	"chatsmith-mini/configs"
	"chatsmith-mini/internal/domain/user"
	"chatsmith-mini/pkg/databases"
	"chatsmith-mini/pkg/oauth"
)

type Module struct {
	Service *Service
}

func NewModule(db databases.PostgresDB, users *user.Service, providers []oauth.Provider, cfg configs.AuthConfig) *Module {
	return &Module{Service: newService(newRepo(db), users, providers, cfg)}
}
