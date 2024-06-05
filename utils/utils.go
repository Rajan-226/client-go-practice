package utils

import (
	"context"
	"fmt"
	"strings"

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

	fmt.Println("Successfully Printed all pods!!")
}

func EditDeploymentImageTag(ctx context.Context, clientSet *kubernetes.Clientset, deploymentName string, tag string) {
	deployment, err := clientSet.AppsV1().Deployments("ns-one").Get(ctx, deploymentName, metav1.GetOptions{})

	if err != nil {
		fmt.Println(err.Error(), " while getting the pods from clientset")
	}

	containers := deployment.Spec.Template.Spec.Containers
	for i := range containers {
		containers[i].Image = updateImageTag(containers[i].Image, tag)
	}

	_, err = clientSet.AppsV1().Deployments("ns-one").Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println(err.Error() + " while updating the image tag of deployment")
	}

	fmt.Printf("Successfully Changed Image of %s deployment to %s!!\n", deploymentName, tag)
}

func updateImageTag(image string, newTag string) string {
	if colonIndex := strings.LastIndex(image, ":"); colonIndex != -1 {
		return image[:colonIndex+1] + newTag
	}

	return image + ":" + newTag
}

func Temp() {
	// informers.NewSharedInformerFactory().Core().V1().Pods().Informer().AddEventHandler()
	// cache.ResourceEventHandlerFuncs
}
