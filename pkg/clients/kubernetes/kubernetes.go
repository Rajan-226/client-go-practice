package kubernetes

import (
	"fmt"
	"time"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	typedClient     *kubernetes.Clientset
	dynamicClient   *dynamic.DynamicClient
	discoveryClient *discovery.DiscoveryClient

	sharedInformerFactory informers.SharedInformerFactory
)

func Init() error {
	cfg, err := getKubeConfig()
	if err != nil {
		return err
	}

	if err = initClients(cfg); err != nil {
		return err
	}

	sharedInformerFactory = informers.NewSharedInformerFactory(typedClient, time.Second*30)

	return nil
}

func initClients(cfg *rest.Config) (err error) {
	typedClient, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf(err.Error() + " while creating the client from config")
	}

	if dynamicClient, err = dynamic.NewForConfig(cfg); err != nil {
		return fmt.Errorf(err.Error() + " while creating the dynamic client from config")
	}

	discoveryClient, err = discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return fmt.Errorf(err.Error() + " while creating the discovery client from config")
	}

	return nil
}

func TypedClient() *kubernetes.Clientset {
	if typedClient == nil {
		panic("nil typed client")
	}

	return typedClient
}

func DyanmicClient() *dynamic.DynamicClient {
	if dynamicClient == nil {
		panic("nil dynamic client")
	}

	return dynamicClient
}

func DiscoveryClient() *discovery.DiscoveryClient {
	if discoveryClient == nil {
		panic("nil discovery client")
	}

	return discoveryClient
}

func SharedInformerFactory() informers.SharedInformerFactory{
	if sharedInformerFactory == nil {
		panic("nil shared informer factory")
	}

	return sharedInformerFactory
}

func getKubeConfig() (*rest.Config, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/home/rajan/.kube/config")
	if err != nil {
		fmt.Println(err.Error() + " while building config from flag")

		cfg, err = rest.InClusterConfig()

		if err != nil {
			return nil, fmt.Errorf(err.Error() + " while getting config from inside the cluster")
		}
	}

	return cfg, nil
}
