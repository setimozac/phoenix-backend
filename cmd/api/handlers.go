package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/setimozac/phoenix-backend/internal/types"
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

	_ = app.WriteJSON(w, http.StatusOK, payload)
}

func (app *application) GetAllEnvManagers(w http.ResponseWriter, r *http.Request) {

	services, err := app.DB.AllEnvManagers()
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, errors.New("bad request"))
		return
	}

	_ = app.WriteJSON(w, http.StatusOK, services)
}

func (app *application) UpdateEnvManagers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Println("Invalide HTTP method")
		app.ErrorJSON(w, errors.New("olny http put method is allowed"), http.StatusMethodNotAllowed)
		return
	}
	var envManagers []types.EnvManager
	var failedToUpdate []int

	err := app.ReadJSON(w, r, &envManagers)
	if err != nil {
		log.Println("Invalide request BODY")
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	log.Println("[]envManagers", envManagers)
	for _, item := range envManagers {
		if app.K8sClient != nil {
			err = app.UpdateEnvManagerInCluster(&item)
			if err != nil {
				log.Println("Failed updating EnvManager:", item)
				failedToUpdate = append(failedToUpdate, item.ID)
			}
		} else {
			log.Println("just for offline testing - UpdateEnvManagers()", item)
			err = app.DB.UpdateEnvManager(&item)
			if err != nil {
				log.Println("Failed updating EnvManager:", item)
			}
		}
	}

	app.WriteJSON(w, http.StatusAccepted, failedToUpdate)
	
}

