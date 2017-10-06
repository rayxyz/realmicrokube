package grpclb

import (
	"log"
	"runtime/debug"
	"time"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/naming"
	kube "k8s.io/client-go/kubernetes"
)

// kubeResolver resolves service names using Kubernetes endpoints.
type kubeResolver struct {
	k8sClient *kube.Clientset
	namespace string
	watcher   *watcher
}

// NewResolver returns a new Kubernetes resolver.
func NewResolver(client *kube.Clientset, namespace string) *kubeResolver {
	if namespace == "" {
		namespace = "default"
	}
	return &kubeResolver{client, namespace, nil}
}

// Resolve creates a Kubernetes watcher for the named target.
func (r *kubeResolver) Resolve(target string) (naming.Watcher, error) {
	resultChan := make(chan watchResult)
	stopCh := make(chan struct{})
	go until(func() {
		err := r.watch(target, stopCh, resultChan)
		if err != nil {
			grpclog.Printf("kuberesolver: watching ended with error='%v', will reconnect again", err)
		}
	}, time.Second, stopCh)

	r.watcher = &watcher{
		target:    target,
		endpoints: make(map[string]interface{}),
		stopCh:    stopCh,
		result:    resultChan,
	}
	return r.watcher, nil
}

func (r *kubeResolver) watch(target string, stopCh <-chan struct{}, resultCh chan<- watchResult) error {
	log.Println("I am watching...")
	// for {
	// 	select {
	// 	case <-stopCh:
	// 		return nil
	// 	case up, more := <-sw.ResultChan():
	// 		if more {
	// 			resultCh <- watchResult{err: nil, ep: &up}
	// 		} else {
	// 			return nil
	// 		}
	// 	}
	// }
	return nil
}

func until(f func(), period time.Duration, stopCh <-chan struct{}) {
	select {
	case <-stopCh:
		return
	default:
	}
	for {
		func() {
			defer handleCrash()
			f()
		}()
		select {
		case <-stopCh:
			return
		case <-time.After(period):
		}
	}
}

// HandleCrash simply catches a crash and logs an error. Meant to be called via defer.
func handleCrash() {
	if r := recover(); r != nil {
		callers := string(debug.Stack())
		grpclog.Printf("kuberesolver: recovered from panic: %#v (%v)\n%v", r, r, callers)
	}
}
