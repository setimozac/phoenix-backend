package main

import (

	"github.com/setimozac/phoenix-backend/internal/repository"
)

const port = 8080

type application struct {
	Domain string
	DSN string
	DB repository.DataBaseRepo
}