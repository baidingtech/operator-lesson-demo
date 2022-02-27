package main

import (
	"fmt"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	//create config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	//create client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	//get informer
	//factory := informers.NewSharedInformerFactory(clientset, 0)
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 0, informers.WithNamespace("default"))
	informer := factory.Core().V1().Pods().Informer()

	//add event handler
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("ADD Event")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("Update Event")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Delete Event")
		},
	})

	//start informer
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)
	<- stopCh

}
