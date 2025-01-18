package dbrepo

import (
	"database/sql"
	"time"
	"github.com/setimozac/phoenix-backend/internal/types"
)

type PgDBRepo struct{
	DBConn *sql.DB
}

const dbTimeout = time.Second * 3

func (pg *PgDBRepo) Connection() interface{} {
	return pg.DBConn
}

func (pg *PgDBRepo) AllEnvManagers() ([]*types.Service, error){
	return nil,nil
}
func (pg *PgDBRepo) GetEnvManagerByName(name string) (*types.Service, error) {
	return nil,nil
}
func (pg *PgDBRepo) GetEnvManagerByIdFromDB(id int) (*types.Service, error) {
	return nil,nil
}