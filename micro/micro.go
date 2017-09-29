package micro

import (
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	kbapiv1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var clientset *kubernetes.Clientset

func init() {
	log.Println("Initializing micro...")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func initKubeOutofCluster() {
	// Out of cluster
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Println("Create client set error.")
		panic(err.Error())
	}
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

type Service struct {
	Config       *ServiceConfig
	NewClientRef interface{}
}

// type KubeService struct {
// 	Host string
// 	Port int
// }

func NewService(config *ServiceConfig, server interface{}, grpcRegisterServer interface{}) {
	listeningAddr := config.Host + ":" + strconv.Itoa(config.Port)
	listener, err := net.Listen("tcp", listeningAddr)
	log.Println("Service listening at =>", listeningAddr)
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}

	grpcServer := grpc.NewServer()
	var args []reflect.Value
	args = append(args, reflect.ValueOf(grpcServer))
	args = append(args, reflect.ValueOf(server))
	reflect.ValueOf(grpcRegisterServer).Call(args)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}
}

func NewServiceClient(service string, newClientRef interface{}) (*Service, error) {
	if service == "" || newClientRef == nil {
		return nil, errors.New("Create service client error. Arguments nil.")
	}
	srv, err := queryKubeService("default", service)
	srvConf := &ServiceConfig{
		Host: srv.Spec.ClusterIP,
		Port: int(srv.Spec.Ports[0].Port),
	}
	if err != nil {
		return nil, err
	}
	return &Service{Config: srvConf, NewClientRef: newClientRef}, nil
}

func queryKubeService(namespace, service string) (*kbapiv1.Service, error) {
	srv, err := clientset.CoreV1().Services(namespace).Get(service, meta.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return srv, nil
}

func (s *Service) Call(method string, ctx context.Context, reqObj interface{}) (interface{}, error) {
	address := s.Config.Host + ":" + strconv.Itoa(s.Config.Port)
	conn, err := grpc.Dial(address)
	if err != nil {
		log.Println("Connection to server error.")
		return nil, err
	}
	if conn == nil {
		return nil, errors.New("Connection cannot be established.")
	}
	defer conn.Close()

	var client reflect.Value
	var newClientArgs []reflect.Value
	newClientArgs = append(newClientArgs, reflect.ValueOf(conn))
	newClientVals := reflect.ValueOf(s.NewClientRef).Call(newClientArgs)
	if newClientVals != nil && len(newClientVals) > 0 {
		client = newClientVals[0]
	}

	if client.IsNil() {
		return nil, errors.New("Parse grpc client error.")
	}

	var methodArgs []reflect.Value
	methodArgs = append(methodArgs, reflect.ValueOf(ctx))
	methodArgs = append(methodArgs, reflect.ValueOf(reqObj))
	// Call grpc method
	methodVals := client.MethodByName(method).Call(methodArgs)

	var respResult interface{}
	var respError error
	if methodVals != nil && len(methodVals) > 0 {
		if methodVals[0].CanInterface() {
			if methodVals[0].Interface() != nil {
				respResult = methodVals[0].Interface()
			}
		}
	}
	if methodVals != nil && len(methodVals) > 1 {
		if methodVals[1].CanInterface() {
			if methodVals[1].Interface() != nil {
				respError = methodVals[1].Interface().(error)
			}
		}
	}

	log.Printf("RespResult => %#v RespError => %#v", respResult, respError)

	return respResult, respError
}
