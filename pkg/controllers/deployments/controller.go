package deployments

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rajan-226/client-go-practice/pkg/clients/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	//client set can be used as a argument to struct to remove dependency on internal kubernetes module
	ctx         context.Context
	cacheSynced cache.InformerSynced
	queue       workqueue.RateLimitingInterface
}

func NewController(ctx context.Context) (*Controller, error) {
	informer := kubernetes.SharedInformerFactory().Apps().V1().Deployments().Informer()

	ctrl := &Controller{
		ctx:         ctx,
		cacheSynced: informer.HasSynced,
		queue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "upips"),
	}

	_, err := informer.AddEventHandler(ctrl)

	if err != nil {
		return nil, fmt.Errorf(err.Error() + " while running deployment informer")
	}

	return ctrl, nil
}

func (c *Controller) Run(stopCh chan struct{}) error {
	go kubernetes.SharedInformerFactory().Apps().V1().Deployments().Informer().Run(stopCh)

	startTime := time.Now()
	if !cache.WaitForCacheSync(stopCh, c.cacheSynced) {
		return fmt.Errorf("cache didn't sync")
	}
	fmt.Printf("Time spent while waiting for cache sync: %d\n", time.Since(startTime).Microseconds())

	// Keep on polling from queue
	// If the queue is shut down, informer and polling from queue both will get stopped
	go func() {
		for {
			item, shutdown := c.queue.Get()
			if shutdown {
				close(stopCh)
				return
			}

			if err := c.process(item); err != nil {
				fmt.Println("Not able to process item due to: " + err.Error())
			} else {
				c.queue.Done(item)
			}
		}
	}()

	return nil
}

func (c *Controller) process(item interface{}) (err error) {
	dep, ok := item.(*appsv1.Deployment)
	if !ok {
		return errors.New("not able to unmarshal item into deployment")
	}

	fmt.Printf("Starting processing for deployment %s in the %s namespace\n\n", dep.GetName(), dep.GetNamespace())

	if isDeleted, err := c.handleIfDeleted(dep); isDeleted {
		return nil
	} else if err != nil {
		return err
	}

	var svc *corev1.Service
	if svc, err = c.createService(dep); err != nil {
		return err
	}

	if _, err := c.createIngress(svc); err != nil {
		return err
	}

	return nil
}

func (c *Controller) OnAdd(obj interface{}, isInInitialList bool) {
	if !isInInitialList {
		c.queue.Add(obj)
	}
}

func (c *Controller) OnDelete(obj interface{}) {
	c.queue.Add(obj)
}

func (c *Controller) createService(dep *appsv1.Deployment) (svc *corev1.Service, err error) {
	clientSet := kubernetes.TypedClient()

	svc = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dep.GetName(),
			Namespace: dep.GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			Selector: dep.Spec.Selector.MatchLabels,
			Ports: []corev1.ServicePort{
				{
					Port: dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort,
				},
			},
		},
	}
	if svc, err = clientSet.CoreV1().Services(dep.GetNamespace()).Create(c.ctx, svc, metav1.CreateOptions{}); err != nil {
		return nil, err
	}
	fmt.Printf("Successfully created service with name %s!!\n\n", dep.GetName())

	return svc, nil
}

func (c *Controller) createIngress(svc *corev1.Service) (ingress *netv1.Ingress, err error) {
	clientSet := kubernetes.TypedClient()

	pathType := "Prefix"
	ingress = &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				{
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     fmt.Sprintf("/%s", svc.Name),
									PathType: (*netv1.PathType)(&pathType),
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: svc.Name,
											Port: netv1.ServiceBackendPort{
												Number: int32(svc.Spec.Ports[0].TargetPort.IntValue()),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if ingress, err = clientSet.NetworkingV1().Ingresses(svc.Namespace).Create(c.ctx, ingress, metav1.CreateOptions{}); err != nil {
		return nil, err
	}

	fmt.Printf("Successfully created ingress with name %s!!\n\n", svc.GetName())

	return ingress, nil
}

func (c *Controller) handleIfDeleted(dep *appsv1.Deployment) (deleted bool, err error) {
	clientset := kubernetes.TypedClient()

	if _, err = clientset.AppsV1().Deployments(dep.GetNamespace()).Get(c.ctx, dep.GetName(), metav1.GetOptions{}); err == nil {
		return false, nil
	}

	if !apierrors.IsNotFound(err) {
		return false, err
	}

	// delete service
	err = clientset.CoreV1().Services(dep.GetNamespace()).Delete(c.ctx, dep.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}
	fmt.Printf("Deleted service %s successfully!!\n\n", dep.GetName())

	//delete ingress
	err = clientset.NetworkingV1().Ingresses(dep.GetNamespace()).Delete(c.ctx, dep.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return false, err
	}
	fmt.Printf("Deleted ingress %s successfully!!\n\n", dep.GetName())


	return true, nil

}

func (c *Controller) OnUpdate(oldObj, newObj interface{}) {
}
