package main

import (
	"encoding/json"
	"fmt"
	// "github.com/setimozac/phoenix-backend/internal/types"
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

	services, err := app.DB.AllEnvManagers()

	data, err := json.Marshal(services)
	if err != nil{
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}