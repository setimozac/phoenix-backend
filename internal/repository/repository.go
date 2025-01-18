package repository

import (
	// "database/sql"

	"github.com/setimozac/phoenix-backend/internal/types"
)

type DataBaseRepo interface {
	Connection() interface{}
	AllEnvManagers() ([]*types.Service, error)
	GetEnvManagerByName(name string) (*types.Service, error)
	GetEnvManagerByIdFromDB(id int) (*types.Service, error)
}

// type DataBaseRepoDynamoDB interface {
// 	// Connection() *sql.DB
// 	AllEnvManagers() ([]*types.Service, error)
// 	GetEnvManagerByName(name string) (*types.Service, error)
// 	GetEnvManagerByIdFromDB(id int) (*types.Service, error)
// }