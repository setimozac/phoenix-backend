package dbrepo

import (
	"context"
	"database/sql"
	"log"
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

func (pg *PgDBRepo) AllEnvManagers() ([]*types.EnvManager, error){
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var envManagers []*types.EnvManager

	query := `
		SELECT id, name, min_replicas, enabled, ui_enabled, last_update FROM env_managers
	`
	rows, err := pg.DBConn.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var envManager types.EnvManager
		err := rows.Scan(
			&envManager.ID,
			&envManager.Name,
			&envManager.MinReplicas,
			&envManager.Enabled,
			&envManager.UIEnabled,
			&envManager.LastUpdate,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		envManagers = append(envManagers, &envManager)
	}

	
	return envManagers,nil
}
func (pg *PgDBRepo) GetEnvManagerByName(name string) (*types.EnvManager, error) {
	var envManager *types.EnvManager
	return envManager,nil
}
func (pg *PgDBRepo) GetEnvManagerByIdFromDB(id int) (*types.EnvManager, error) {
	return nil,nil
}

func (pg *PgDBRepo) InsertEnvManager(em types.EnvManager) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `INSERT INTO env_managers(name, min_replicas, enabled, last_update) VALUES($1,$2,$3,$4)`
	em.LastUpdate = time.Now().Unix()

	err := pg.DBConn.QueryRowContext(ctx, stmt, em.Name, em.MinReplicas, em.Enabled, em.LastUpdate).Scan(&newID)
	if err != nil {
		return 0, nil
	}
	return newID, nil
}