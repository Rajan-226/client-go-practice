package main

import (
	"context"
	"flag"
	"fmt"

	"client-go-practice/utils"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	ctx := context.Background()
	clientSet, err := setupClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = utils.PrintPods(ctx, clientSet); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = utils.EditDeploymentImageTag(ctx, clientSet, "nginx", "1.21.6"); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func setupClient() (*kubernetes.Clientset, error) {
	kubeConfig := flag.String("kubeconfig", "/home/rajan/.kube/aconfig", "location to kube config")
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		fmt.Println(err, " while building config from flag")

		cfg, err = rest.InClusterConfig()

		if err != nil {
			return nil, fmt.Errorf(err.Error() + " while getting config from inside the cluster")
		}
	}

	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf(err.Error() + " while creating the client from config")
	}

	return clientSet, nil
}
