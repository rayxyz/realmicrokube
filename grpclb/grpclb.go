package grpclb

import (
	"log"

	kbapiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientset *kubernetes.Clientset

type loadBalancer struct {
}

func init() {
	initKubeInCluster()
}

func initKubeInCluster() {
	// In cluster
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func (lb *loadBalancer) queryService(namespace, service string) (*kbapiv1.Service, error) {
	srv, err := clientset.CoreV1().Services(namespace).Get(service, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Got the service => ", service, " in gRPC load balancer.")
	return srv, nil
}
