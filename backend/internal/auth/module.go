package auth

import (
	"chatsmith-mini/configs"
	"chatsmith-mini/internal/user"
)

type CreateModuleDTO struct {
	Config configs.AuthConfig
	Users  *user.Service
}

type Module struct {
	Service *Service
}

func NewModule(dto *CreateModuleDTO) *Module {
	return &Module{Service: &Service{
		google: newGoogle(dto.Config),
		codec:  newCodec(dto.Config.SessionSecret),
		users:  dto.Users,
		cfg:    dto.Config,
	}}
}
