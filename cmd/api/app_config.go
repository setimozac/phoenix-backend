package main

import (

	"github.com/setimozac/phoenix-backend/internal/repository"
	"k8s.io/client-go/dynamic"
)

const port = 8080

type application struct {
	Domain string
	DSN string
	DB repository.DataBaseRepo
	K8sActive bool
	K8sClient *dynamic.DynamicClient
}