package main

import (
	"flag"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	_ = setupClient()

	
}

func setupClient() *kubernetes.Clientset {
	kubeConfig := flag.String("kubeconfig", "/home/rajan/.kube/config", "location to kube config")
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		throwError(err, "building config from flag")
	}

	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		throwError(err, "creating the client from config")
	}

	return clientSet
}

func throwError(err error, when string) {
	fmt.Println(err.Error() + " while " + when)
}
