package dbrepo

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgconn"
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
			&envManager.MinReplica,
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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var envManager types.EnvManager
	query := `
		SELECT id, name, min_replicas, enabled, ui_enabled, last_update, namespace, cr_name
		FROM env_managers
		WHERE name = $1
		`

	row := pg.DBConn.QueryRowContext(ctx, query, name)
	err := row.Scan(
		&envManager.ID,
		&envManager.Name,
		&envManager.MinReplica,
		&envManager.Enabled,
		&envManager.UIEnabled,
		&envManager.LastUpdate,
		&envManager.Metadata.Namespace,
		&envManager.Metadata.Name,
	)
	if err != nil{
		return nil, err
	}
	return &envManager, nil
}

func (pg *PgDBRepo) GetEnvManagerById(id int) (*types.EnvManager, error) {
	return nil,nil
}

func (pg *PgDBRepo) UpdateEnvManager(em *types.EnvManager) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `
		UPDATE env_managers
		SET min_replicas = $1, enabled = $2
		WHERE name = $3;
	`

	_, err := pg.DBConn.ExecContext(ctx, stmt, em.MinReplica, em.Enabled, em.Name)
	if err != nil{
		return err
	}

	return nil
}

func (pg *PgDBRepo) InsertEnvManager(em *types.EnvManager) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `INSERT INTO env_managers(name, min_replicas, enabled, last_update) VALUES($1,$2,$3,$4,$5,$6) RETURNING id;`
	em.LastUpdate = time.Now().Unix()

	err := pg.DBConn.QueryRowContext(ctx, stmt, em.Name, em.MinReplica, em.Enabled, em.LastUpdate, em.Metadata.Namespace, em.Metadata.Name).Scan(&newID)
	log.Println("id:", newID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				log.Println("duplicate key(name) error for:", em.Name)
			}
		}
		return 0, err
	}
	return newID, nil
}

func (pg *PgDBRepo) DelteEnvManager(em *types.EnvManager) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		DELETE FROM env_managers WHERE name = $1;
	`

	_, err := pg.DBConn.ExecContext(ctx, stmt, em.Name)
	if err != nil{
		log.Println("unable to delete a record. envManager: ", em.Name)
		return err
	}

	return nil
}