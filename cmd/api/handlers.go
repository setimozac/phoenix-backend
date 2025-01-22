package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/setimozac/phoenix-backend/internal/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

func (app *application) UpdateEnvManagers(w http.ResponseWriter, r *http.Request) {
	if app.K8sClient != nil {

	} else {
		
	}
}

func (app *application) UpdateEnvManager(em *types.EnvManager) error {
	namespace := em.Metadata.Namespace
	crName := em.Metadata.Name

	currentCR, err := app.K8sClient.Resource(app.GRV).Namespace(namespace).Get(context.TODO(), crName,metav1.GetOptions{})
	if err != nil {
		log.Println("unable to get the CR", crName, err)
		return err
	}

	spec, found, err := unstructured.NestedMap(currentCR.Object, "spec")
	if err != nil {
		log.Println("unable to get the spec",  err)
		return err
	}
	if !found {
		log.Println("spec not found", crName)
		return nil
	}

	spec["enabled"] = em.Enabled
	spec["lastUpdate"] = em.LastUpdate
	spec["minReplica"] = em.MinReplica
	spec["uiEnabled"] = em.UIEnabled
	
	// update the current CR with new spec
	err = unstructured.SetNestedField(currentCR.Object, spec, "spec")
	if err != nil {
		log.Println("unable to update the spec", err)
		return err
	}

	updatedCR, err := app.K8sClient.Resource(app.GRV).Namespace(namespace).Update(context.TODO(), currentCR, metav1.UpdateOptions{})
	if err != nil {
		log.Println("unable to update the CR:",em.Metadata.Name, err)
		return err
	}
	log.Println("cr was updated successfully", updatedCR.GetName())
	return nil
}