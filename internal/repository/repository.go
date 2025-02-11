package repository

import (
	// "database/sql"

	"github.com/setimozac/phoenix-backend/internal/types"
)

type DataBaseRepo interface {
	Connection() interface{}
	AllEnvManagers() ([]*types.EnvManager, error)
	GetEnvManagerByName(name string) (*types.EnvManager, error)
	GetEnvManagerById(id int) (*types.EnvManager, error)
	InsertEnvManager(em *types.EnvManager) (int, error)
	UpdateEnvManager(em *types.EnvManager) error
	DelteEnvManager(em *types.EnvManager) error
	AddEvents(events []string, servceName string) error
}