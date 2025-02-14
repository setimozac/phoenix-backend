package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
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
		SELECT em.id, em.name, em.min_replicas, em.enabled, em.ui_enabled, em.last_update, em.namespace, em.cr_name, COALESCE(array_agg(e.event_name), '{}') as events
		FROM env_managers em
		LEFT JOIN events e ON em.name = e.service_name
		GROUP BY em.id, em.name, em.min_replicas, em.enabled, em.ui_enabled, em.last_update, em.namespace, em.cr_name;
	`
	rows, err := pg.DBConn.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var namespace, crName string
	var events pgtype.TextArray
	for rows.Next() {
		var envManager types.EnvManager
		err := rows.Scan(
			&envManager.ID,
			&envManager.Name,
			&envManager.MinReplica,
			&envManager.Enabled,
			&envManager.UIEnabled,
			&envManager.LastUpdate,
			&namespace,
			&crName,
			&events,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		envManager.Metadata = &types.Metadata{
			Name: crName,
			Namespace: namespace,
		}
		if events.Status == pgtype.Present {
			for _, elem := range events.Elements {
				envManager.Events = append(envManager.Events, elem.String)
			}
		}
		

		envManagers = append(envManagers, &envManager)
	}

	
	return envManagers,nil
}
func (pg *PgDBRepo) GetEnvManagerByName(name string) (*types.EnvManager, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	log.Println("From GetEnvManagerByName", name)
	var envManager types.EnvManager
	var metaNamespace string
	var metaName string
	// query := `
	// 	SELECT id, name, min_replicas, enabled, ui_enabled, last_update, namespace, cr_name
	// 	FROM env_managers
	// 	WHERE name = $1
	// 	`
	query := `
		SELECT em.id, em.name, em.min_replicas, em.enabled, em.ui_enabled, em.last_update, em.namespace, em.cr_name, COALESCE(array_agg(e.event_name), '{}') as events
		FROM env_managers em
		LEFT JOIN events e ON em.name = e.service_name
		WHERE em.name = $1
		GROUP BY em.id, em.name, em.min_replicas, em.enabled, em.ui_enabled, em.last_update, em.namespace, em.cr_name;
		`
		var events pgtype.TextArray
	row := pg.DBConn.QueryRowContext(ctx, query, name)
	err := row.Scan(
		&envManager.ID,
		&envManager.Name,
		&envManager.MinReplica,
		&envManager.Enabled,
		&envManager.UIEnabled,
		&envManager.LastUpdate,
		&metaNamespace,
		&metaName,
		&events,
	)
	if err != nil{
		return nil, err
	}
	envManager.Metadata = &types.Metadata{
		Name: metaName,
		Namespace: metaNamespace,
	}
	if events.Status == pgtype.Present {
		for _, elem := range events.Elements {
			envManager.Events = append(envManager.Events, elem.String)
		}
	}
	return &envManager, nil
}

func (pg *PgDBRepo) GetEnvManagerById(id int) (*types.EnvManager, error) {
	return nil,nil
}

func (pg *PgDBRepo) UpdateEnvManager(em *types.EnvManager) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	em.LastUpdate = time.Now().Unix()
	stmt := `
		UPDATE env_managers
		SET min_replicas = $1, enabled = $2, ui_enabled=$3, last_update=$4
		WHERE name = $5;
	`

	_, err := pg.DBConn.ExecContext(ctx, stmt, em.MinReplica, em.Enabled,em.UIEnabled, em.LastUpdate, em.Name)
	if err != nil{
		return err
	}

	err = pg.UpdateEvents(em.Events, em.Name)
	if err != nil{
		return err
	}

	return nil
}

func (pg *PgDBRepo) InsertEnvManager(em *types.EnvManager) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `INSERT INTO env_managers(name, min_replicas, enabled, last_update, namespace, cr_name) VALUES($1,$2,$3,$4,$5,$6) RETURNING id;`
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
	if len(em.Events) > 0 {
		err = pg.AddEvents(em.Events, em.Name)
		if err != nil {
			log.Println(err)
			return newID, err
		}
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

func (pg *PgDBRepo) AddEvents(events []string, servceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := "INSERT INTO events (event_name, service_name) Values "
	var values []interface{}

	for i, event := range events {
		stmt += fmt.Sprintf("($%d, $%d),", i*2+1, i*2+2)
		values = append(values, event, servceName)
	}

	stmt = stmt[:len(stmt)-1]        		 // Remove last comma
	log.Println(stmt)
	log.Println(events)
	log.Println(values)
	_, err := pg.DBConn.ExecContext(ctx, stmt, values...)
	if err != nil{
		log.Println("unable to insert events for: ", servceName, "Err: ", err)
		return err
	}

	return nil
}

func (pg *PgDBRepo) UpdateEvents(events []string, servceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	tx, err := pg.DBConn.BeginTx(ctx, nil)
	if err != nil {	
		return err
	}

	stmt := `DELETE FROM events WHERE service_name = $1`

	_, err = tx.Exec(stmt, servceName)
	if err != nil {
		log.Panicln("failed to delete old events", err)
		tx.Rollback()
		return err
	}
	tx.Commit()

	err = pg.AddEvents(events, servceName)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PgDBRepo) GetAllEvents() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	
	var events []string
	query := `SELECT DISTINCT event_name FROM events`

	rows, err := pg.DBConn.QueryContext(ctx, query)
	if err != nil {
		return nil,err
	}
	defer rows.Close()

	for rows.Next() {
		var eventName string
		if err := rows.Scan(&eventName); err != nil {
			return nil,err
		}
		events = append(events, eventName)
	}

	return events,nil
}