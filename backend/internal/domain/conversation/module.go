package conversation

import "chatsmith-mini/pkg/databases"

type Module struct {
	Service *Service
}

func NewModule(db databases.PostgresDB) *Module {
	return &Module{Service: newService(newRepo(db))}
}
