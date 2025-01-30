package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/setimozac/phoenix-backend/internal/repository/dbrepo"
	"github.com/setimozac/phoenix-backend/internal/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var (
	gvr = schema.GroupVersionResource{
		Group: "phoenix.setimozak",
		Version: "v1beta1",
		Resource: "envmanagers",
	}
)

type Spec struct {
	Enabled bool `json:"enabled"`
	MinReplica int `json:"minReplica"`
	UIEnabled bool `json:"uiEnabled,omitempty"`
	LastUpdate int64 `json:"lastUpdate,omitempty"`
	Name string `json:"name"`
}

func main() {
	// application config
	var app application
	app.GRV = gvr


	// read args from the cli
	flag.StringVar(&app.DSN, "dsn", "host=postgres port=5432 user=postgres password=postgres dbname=env_manager timezone=UTC connect_timeout=5", "db connection string")
	flag.BoolVar(&app.K8sActive, "is-cluster-ready", true, "set it to false if you want to test the endpoints without a cluster running on the server")
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


	// cluster congiguration
	if app.K8sActive {
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Println("error creating in-cluster config.", err)
		}
		dynamicClient, err := dynamic.NewForConfig(config)
		if err != nil {
			log.Println("error creating dynamic client.", err)
		}

		app.K8sClient = dynamicClient

	} else {
		app.K8sClient = nil
	}
	
	// watch CRs if the dynamic client is available
	if app.K8sClient != nil {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return app.K8sClient.Resource(gvr).Namespace(metav1.NamespaceAll).List(context.TODO(), options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return app.K8sClient.Resource(gvr).Namespace(metav1.NamespaceAll).Watch(context.TODO(), options)
				},
			},
			&unstructured.Unstructured{},
			time.Minute*1,
			cache.Indexers{
				cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
			},
		)

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				envManager := obj.(*unstructured.Unstructured)
				log.Println("<=========== CR was added... ===========>")

				var em types.EnvManager
				err := runtime.DefaultUnstructuredConverter.FromUnstructured(envManager.Object["spec"].(map[string]interface{}), &em)
				if err != nil {
					log.Println("unable to extract spec from the object.", err)
				}

				var metadata types.Metadata
				err = runtime.DefaultUnstructuredConverter.FromUnstructured(envManager.Object["metadata"].(map[string]interface{}), &metadata)
				if err != nil {
					log.Println("unable to extract metadata from the object.", err)
				}

				em.Metadata = &metadata
				log.Println("em with metadata: ", em.Metadata)


				id, err := app.DB.InsertEnvManager(&em)
				if err != nil {
					log.Println("unable to insert envmanager into the db", err)
					return
				}
				log.Println("obj inserted to the DB with ID: ", id)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				// oldEnvManager := oldObj.(*unstructured.Unstructured)
				newEnvManager := newObj.(*unstructured.Unstructured)
				log.Println("<=========== CR was updated: ===========>", newEnvManager)

				var em types.EnvManager
				err := runtime.DefaultUnstructuredConverter.FromUnstructured(newEnvManager.Object["spec"].(map[string]interface{}), &em)
				if err != nil {
					log.Println("unable to extract spec from the object.", err)
				}
				log.Println("em", em)
				log.Println("em.Name", em.Name)
				oldEm, err := app.DB.GetEnvManagerByName(em.Name)
				if err != nil {
					log.Println("unable to fetch envmanager from the db", err)
					return
				}
				oldEm.Enabled = em.Enabled
				oldEm.MinReplica = em.MinReplica

				log.Println("em", em)
				log.Println("oldEm", oldEm)


				err = app.DB.UpdateEnvManager(oldEm)
				if err != nil {
					log.Println("unable to update the envmanager",oldEm.Name, err)
					return
				}
				log.Println("obj was updated successfully. envmanager.name: ",oldEm.Name)
				
			},
			DeleteFunc: func(obj interface{}) {
				var em types.EnvManager
				envManager, ok := obj.(*unstructured.Unstructured); if !ok{
					tombStone, ok := obj.(cache.DeletedFinalStateUnknown); if !ok {
						log.Println("failed to get the object from tombstone")
						return
					}
					envManager, ok := tombStone.Obj.(*unstructured.Unstructured); if !ok {
						log.Println("tombstone contained object that is not unstructured")
						return
					}
					err := runtime.DefaultUnstructuredConverter.FromUnstructured(envManager.Object["spec"].(map[string]interface{}), &em)
					if err != nil {
						log.Println("unable to extract spec from the object.", err)
					}
					app.DB.DelteEnvManager(&em)
					log.Println("CR was deleted: ", em)
				}
				
				err := runtime.DefaultUnstructuredConverter.FromUnstructured(envManager.Object["spec"].(map[string]interface{}), &em)
				if err != nil {
					log.Println("unable to extract spec from the object.", err)
				}
				app.DB.DelteEnvManager(&em)
				log.Println("CR was deleted: ", em)
			},
		})

		stopCh := make(chan struct{})
		defer close(stopCh)
		log.Println(">>>>>>> Starting the k8s informer <<<<<<<<")
		go informer.Run(stopCh)

	}
	


	


	// start the server
	app.Domain = "backend.phoenix.com"

	log.Println("Starting the server on port:", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}