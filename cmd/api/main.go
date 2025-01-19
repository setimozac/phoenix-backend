package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/setimozac/phoenix-backend/internal/repository/dbrepo"
)

func main() {
	// application config
	var app application

	// read args from the cli
	flag.StringVar(&app.DSN, "dsn", "host=postgres port=5432 user=postgres password=postgres dbname=env_manager timezone=UTC connect_timeout=5", "db connection string")

	flag.Parse()

	// connect to DB
	conn, err := app.connectToPGDB()
	if err != nil {
		log.Fatal("cannot connect to the db", err)
	}

	app.DB = &dbrepo.PgDBRepo{
		DBConn: conn,
	}
	defer app.DB.Connection().(*sql.DB).Close()

	// start the server
	app.Domain = "backend.phoenix.com"

	log.Println("Starting the server on port:", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}