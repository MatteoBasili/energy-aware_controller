package main

import (
	"log"

	"ea-controller/internal/config"
	"ea-controller/internal/controller"
	"ea-controller/internal/informer"
	"ea-controller/internal/store"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	istioclient "istio.io/client-go/pkg/clientset/versioned"
)

func main() {
	cfg := config.Load()

	restConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	istioClient, err := istioclient.NewForConfig(restConfig)
	if err != nil {
		log.Fatal(err)
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		log.Fatal(err)
	}

	store := store.New()

	informer.Start(dynamicClient, cfg.Namespace, store)

	go controller.Run(cfg, store, istioClient)

	select {}
}
