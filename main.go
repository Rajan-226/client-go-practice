package main

import (
	"context"
	"flag"
	"fmt"

	"client-go-practice/utils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	ctx := context.Background()
	clientSet := setupClient()

	// utils.PrintPods(ctx, clientSet)
	utils.EditDeploymentImageTag(ctx, clientSet, "nginx", "1.21.6")
}

func setupClient() *kubernetes.Clientset {
	kubeConfig := flag.String("kubeconfig", "/home/rajan/.kube/config", "location to kube config")
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		fmt.Println(err, " while building config from flag")
	}

	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		fmt.Println(err, " while creating the client from config")
	}

	return clientSet
}
