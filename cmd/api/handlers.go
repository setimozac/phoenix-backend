package main

import (
	"encoding/json"
	"fmt"
	"github.com/setimozac/phoenix-backend/internal/types"
	"net/http"
)

func (app *application) Hello(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Status string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status: "active",
		Message: "Phoenix backend is up",
		Version: "0.0.1",
	}

	data, err := json.Marshal(payload)
	if err != nil{
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (app *application) GetAllEnvManagers(w http.ResponseWriter, r *http.Request) {
	var services []types.Service

	service1 := types.Service{
		ID: 1,
		Enable: true,
		MinReplicas: 0,
		Name: "service1",
	}
	service2 := types.Service{
		ID: 2,
		Enable: true,
		MinReplicas: 0,
		Name: "service2",
	}
	service3 := types.Service{
		ID: 3,
		Enable: false,
		MinReplicas: 1,
		Name: "service3",
	}
	service4 := types.Service{
		ID: 4,
		Enable: true,
		MinReplicas: 1,
		Name: "service4",
	}

	services = append(services, service1)
	services = append(services, service2)
	services = append(services, service3)
	services = append(services, service4)

	data, err := json.Marshal(services)
	if err != nil{
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}