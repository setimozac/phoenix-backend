package main

import (
	"context"
	"log"
	
	"github.com/setimozac/phoenix-backend/internal/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (app *application) UpdateEnvManagerInCluster(em *types.EnvManager) error {
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