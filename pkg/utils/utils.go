package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/Rajan-226/client-go-practice/pkg/clients/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func PrintPods(ctx context.Context) error {
	client := kubernetes.TypedClient()

	pods, err := client.CoreV1().Pods("ns-one").List(ctx, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf(err.Error(), " while getting the pods from client")
	}

	for index, pod := range pods.Items {
		fmt.Printf("Pod %d : %+v\n", index, pod.Name)
	}

	fmt.Printf("Successfully Printed all pods!!\n\n")
	return nil
}

func EditDeploymentImageTag(ctx context.Context, deploymentName string, tag string) error {
	client := kubernetes.TypedClient()

	deployment, err := client.AppsV1().Deployments("ns-one").Get(ctx, deploymentName, metav1.GetOptions{})

	if err != nil {
		return fmt.Errorf(err.Error() + " while getting the pods from client")
	}

	containers := deployment.Spec.Template.Spec.Containers
	for i := range containers {
		containers[i].Image = updateImageTag(containers[i].Image, tag)
	}

	_, err = client.AppsV1().Deployments("ns-one").Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf(err.Error() + " while updating the image tag of deployment")
	}

	fmt.Printf("Successfully Changed Image of %s deployment to %s!!\n\n", deploymentName, tag)
	return nil
}

func ListResources(ctx context.Context, namespaceName string, resourceName string) error {
	client := kubernetes.DyanmicClient()

	gvr, err := findGVR(ctx, resourceName)
	if err!=nil{
		return fmt.Errorf(err.Error() + " while finding GVR from resource name using discovery client")
	}

	resources, err := client.Resource(gvr).Namespace(namespaceName).List(ctx, metav1.ListOptions{})

	if err != nil {
		return fmt.Errorf(err.Error() + " while listing resources by dynamic client")
	}

	for index, resource := range resources.Items {
		fmt.Printf("%s %d : %+v\n", resource.GetKind(), index, resource.GetName())

	}
	
	fmt.Printf("Successfully Printed all resources!!\n\n")

	return nil
}

func updateImageTag(image string, newTag string) string {
	if colonIndex := strings.LastIndex(image, ":"); colonIndex != -1 {
		return image[:colonIndex+1] + newTag
	}

	return image + ":" + newTag
}

func findGVR(ctx context.Context, resource string) (schema.GroupVersionResource, error) {
	resourceList, err := kubernetes.DiscoveryClient().ServerPreferredResources()
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	for _, resourceGroup := range resourceList {
		for _, apiResource := range resourceGroup.APIResources {
			if strings.Contains(apiResource.Name, "/") {
				continue // Skip sub-resources
			}
			if apiResource.Name == resource {
				groupVersion, err := schema.ParseGroupVersion(resourceGroup.GroupVersion)
				if err != nil {
					return schema.GroupVersionResource{}, err
				}
				return schema.GroupVersionResource{
					Group:    groupVersion.Group,
					Version:  groupVersion.Version,
					Resource: apiResource.Name,
				}, nil
			}
		}
	}

	return schema.GroupVersionResource{}, fmt.Errorf("resource type %s not found", resource)
}

func Temp() {
	// informers.NewSharedInformerFactory().Core().V1().Pods().Informer().AddEventHandler()
	// cache.ResourceEventHandlerFuncs
}
