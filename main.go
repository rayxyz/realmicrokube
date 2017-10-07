package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	db "realmicrokube/service/db"

	"realmicrokube/micro"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientset *kubernetes.Clientset

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func init() {
	doInit()
}

func doInit() {
	// Out of cluster
	// var kubeconfig *string
	// if home := homeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }
	// flag.Parse()
	//
	// // use the current context in kubeconfig
	// config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// create the clientset
	// clientset, err = kubernetes.NewForConfig(config)
	// if err != nil {
	// 	log.Println("Create client set error.")
	// 	panic(err.Error())
	// }

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

func checkPods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Checking pods...")
	for {
		pods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		realPods := pods.Items
		for i := 0; i < len(realPods); i++ {
			if realPods[i].GetNamespace() == "default" {
				pod, err := clientset.CoreV1().Pods("default").Get(realPods[i].GetName(), metav1.GetOptions{})
				if errors.IsNotFound(err) {
					fmt.Printf("Pod not found\n")
				} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
					fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
				} else if err != nil {
					panic(err.Error())
				} else {
					fmt.Printf("Found pod\n")
				}
				fmt.Println("pod cluster name => ", pod.GetClusterName(), "pod name => ", pod.GetName(), "pod namespace => ", pod.GetNamespace())
			}
		}

		time.Sleep(5 * time.Second)
	}
}

type PodObj struct {
	Name      string `json:"name"`
	NameSpace string `json:"name_space"`
}

func showPods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Show pods.")
	pods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if pods == nil || pods.Items == nil || len(pods.Items) <= 0 {
		return
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	var podList []*PodObj
	for _, pod := range pods.Items {
		fmt.Println("Pod Name => ", pod.GetName())
		podObj := &PodObj{
			Name:      pod.GetName(),
			NameSpace: pod.GetNamespace(),
		}
		podList = append(podList, podObj)
	}
	podListMarshaled, _ := json.Marshal(podList)
	w.Write(podListMarshaled)
	w.Write([]byte(fmt.Sprintf("There are %d pods in the cluster\n", len(pods.Items))))
}

func int32Ptr(i int32) *int32 { return &i }

func deployService(w http.ResponseWriter, r *http.Request) {
	deployConfig := &micro.KubeServiceDeployConfig{
		Namespace:  apiv1.NamespaceDefault,
		Name:       "com.shendu.service.usercenter.user",
		Port:       int32(83),
		TargetPort: int32(9999),
		Image:      "ray-xyz.com:9090/realmicroserver",
		Replicas:   1,
	}
	success, desc := micro.DeployKubeService(deployConfig)
	if !success {
		w.Write([]byte(fmt.Sprintf("Deploy kubernetes service failed. Deployment description => %s", desc)))
		return
	}
	w.Write([]byte("Deploy kubernetes service => " + deployConfig.Name + ", image => " + deployConfig.Image))
}

func queryUCount(w http.ResponseWriter, r *http.Request) {
	db := db.NewDB()
	count, err := db.QueryUserCount()
	defer func() {
		if db != nil {
			db.Close()
		}
	}()
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println("The user count => ", count)
	w.Write([]byte("The user count is => " + strconv.Itoa(count)))
}

func main() {
	if clientset == nil {
		panic("Client set is nil.")
	}

	port := "7878"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am in real micro kube server."))
	})
	http.HandleFunc("/showpods", showPods)
	http.HandleFunc("/checkpods", checkPods)
	http.HandleFunc("/deploysvc", deployService)
	http.HandleFunc("/ucount", queryUCount)
	log.Println("Server running on port => ", port)
	http.ListenAndServe(":"+port, nil)
}
