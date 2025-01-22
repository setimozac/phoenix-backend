package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "github.com/setimozac/phoenix-backend/internal/types"
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
	if err != nil {
		log.Println(err)
		return
	}

	data, err := json.Marshal(services)
	if err != nil{
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// func (app *application) TestGetEnvManager(w http.ResponseWriter, r *http.Request) {
// 	em := types.EnvManager{
// 		Name: "service1",
// 		MinReplica: 2,
// 		Enabled: true,
// 	}

// 	id, err := app.DB.InsertEnvManager(&em)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	log.Println("ID:", id)

// 	envManager, err := app.DB.GetEnvManagerByName(em.Name)
// 	if err != nil {
// 		log.Println("GetEnvManagerByName(em.Name)",err)
// 		return
// 	}

// 	log.Println(envManager.Name)
// 	envManager.Enabled = false
// 	err = app.DB.UpdateEnvManager(envManager)
// 	if err != nil{
// 		fmt.Println("update",err)
// 	}

// 	updatedEnvManager, err := app.DB.GetEnvManagerByName(envManager.Name)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	data, err := json.Marshal(updatedEnvManager)
// 	if err != nil{
// 		fmt.Println(err)
// 	}

// 	w.Header().Set("Content-Type", "Application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(data)
// }