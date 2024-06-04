package utils

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func PrintPods(ctx context.Context, clientSet *kubernetes.Clientset) {
	pods, err := clientSet.CoreV1().Pods("ns-one").List(ctx, metav1.ListOptions{})

	if err != nil {
		fmt.Println(err.Error(), " while getting the pods from clientset")
	}

	for index, pod := range pods.Items {
		fmt.Printf("Pod %d : %+v\n", index, pod.Name)
	}
}
