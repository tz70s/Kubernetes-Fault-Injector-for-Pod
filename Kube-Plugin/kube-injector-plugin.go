package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func createNew() *v1.Pod {

	var newPod v1.Pod
	newPod.TypeMeta.Kind = "Pod"
	newPod.TypeMeta.APIVersion = "v1"
	newPod.ObjectMeta.Name = "node-server"
	newPod.ObjectMeta.Namespace = "default"

	var podSpec v1.PodSpec

	var conSpec v1.Container

	conSpec.Name = "test"
	conSpec.Image = "tz70s/node-server"

	var conPort v1.ContainerPort
	conPort.ContainerPort = 8080
	conSpec.Ports = []v1.ContainerPort{conPort}
	podSpec.Containers = []v1.Container{conSpec}

	newPod.Spec = podSpec

	return &newPod
}

func main() {

	kubeconfig := flag.String("kubeconfig", "/home/tzuchiao/.kube/config", "absolute path to kube config file")
	flag.Parse()

	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	clientset.CoreV1().Pods("default").Create(createNew())

	for {
		pods, err := clientset.CoreV1().Pods("default").List(v1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		for i := 0; i < len(pods.Items); i++ {
			fmt.Println("Pod Namespace : " + pods.Items[i].Namespace)
			fmt.Println("Pod name : " + pods.Items[i].Name)
			fmt.Println("Pod Labels : ")
			labels := pods.Items[i].Labels
			fmt.Println(labels)
			fmt.Println(pods.Items[i].Spec)
		}
		time.Sleep(5 * time.Second)
	}
}
