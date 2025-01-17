package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// application config
	var app application

	// read args from the cli

	// connect to DB

	// start the server
	app.Domain = "backend.phoenix.com"

	log.Println("Starting the server on port:", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}