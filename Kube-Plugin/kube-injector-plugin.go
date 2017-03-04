package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

// Sample create

func createNew() *v1.Pod {

	var newPod v1.Pod
	newPod.TypeMeta.Kind = "Pod"
	newPod.TypeMeta.APIVersion = "v1"
	newPod.ObjectMeta.Name = "node-server"
	newPod.ObjectMeta.Namespace = "default"
	newPod.ObjectMeta.Labels = map[string]string{
		"injector": "false",
	}
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

// Update Pod with fault-injector container
func addinjector(originalPod *v1.Pod) *v1.Pod {

	var conSpec v1.Container
	conSpec.Name = "injector"
	conSpec.Image = "tz70s/fault-injector"
	var conPort v1.ContainerPort
	conPort.ContainerPort = 8282
	conSpec.Ports = []v1.ContainerPort{conPort}

	var newPod v1.Pod
	newPod.TypeMeta.Kind = "Pod"
	newPod.TypeMeta.APIVersion = "v1"
	newPod.ObjectMeta.Name = originalPod.ObjectMeta.Name + "-inject"
	newPod.ObjectMeta.Namespace = "default"
	newPod.ObjectMeta.Labels = map[string]string{
		"injector": "true",
	}

	newPod.Spec.Containers = append(originalPod.Spec.Containers, conSpec)
	return &newPod
}

func main() {

	// set kube-apiserver config file with address, ca, auth-token, etc.
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

	// create sample pod

	clientset.CoreV1().Pods("default").Create(createNew())

	// polling api server until if there have pods with labels - injector = false
	for {
		pods, err := clientset.CoreV1().Pods("default").List(v1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for i := 0; i < len(pods.Items); i++ {
			if pods.Items[i].Labels["injector"] == "false" {
				tmpPod, _ := clientset.CoreV1().Pods("default").UpdateStatus(&pods.Items[i])
				if tmpPod.Status.Phase == v1.PodRunning {
					clientset.CoreV1().Pods("default").Delete(tmpPod.Name, &v1.DeleteOptions{})
					// wait a little bit for status trans to terminate
					fmt.Println("Add Injector!")
					// async handle create
					go clientset.CoreV1().Pods("default").Create(addinjector(tmpPod))
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}
