package refinery

import (
	"log"

	"github.com/kubernetes-sigs/kubebuilder/pkg/builders"

	"github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	"github.com/lwolf/kubereplay/pkg/client/clientset_generated/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	listers "github.com/lwolf/kubereplay/pkg/client/listers_generated/kubereplay/v1alpha1"
	"github.com/lwolf/kubereplay/pkg/controller/sharedinformers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// Created by "kubebuilder create resource" for you to implement controller logic for the Refinery resource API

// Reconcile handles enqueued messages
func (c *RefineryControllerImpl) Reconcile(u *v1alpha1.Refinery) error {
	log.Printf("Running reconcile Refinery for %s\n", u.Name)

	sClient := c.cs.CoreV1().Services(u.Namespace)
	service := GenerateService(u.Name, &u.Spec)
	svc, _ := sClient.Get(service.Name, metav1.GetOptions{})
	if svc == nil {
		_, err := sClient.Create(service)
		if err != nil {
			log.Printf("Failed to create service: %v", err)
			return err
		}
	} else {
		// TODO: compare deployed version to the new one, and update if needed
		log.Printf("service %s/%s exists", service.Namespace, service.Name)
	}

	dClient := c.cs.AppsV1().Deployments(u.Namespace)
	deployment := GenerateDeployment(u.Name, u)
	d, _ := dClient.Get(deployment.Name, metav1.GetOptions{})
	if d == nil {
		// Create Deployment
		log.Printf("Creating refinery deployment...")
		result, err := dClient.Create(deployment)

		if err != nil {
			log.Printf("Failed to create deployment: %v", err)
			return err
		}
		log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	} else {
		// TODO: compare deployed version to the new one, and update if needed
		log.Printf("deployment %s/%s exists", deployment.Namespace, deployment.Name)
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Refinery before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver

		result, getErr := c.Get(u.Namespace, u.Name)
		if getErr != nil {
			log.Fatalf("Failed to get latest version of Silo: %v", getErr)
		}
		result.Status.Deployed = true
		_, updateErr := c.cset.KubereplayV1alpha1().Refineries(result.Namespace).Update(result)
		return updateErr
	})
	if retryErr != nil {
		log.Printf("Update failed: %v", retryErr)
		return retryErr
	}

	return nil
}

// +controller:group=kubereplay,version=v1alpha1,kind=Refinery,resource=refineries
type RefineryControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about Refinery
	lister listers.RefineryLister

	cset *clientset.Clientset
	cs   *kubernetes.Clientset
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *RefineryControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {

	c.lister = arguments.GetSharedInformers().Factory.Kubereplay().V1alpha1().Refineries().Lister()

	c.cs = arguments.GetSharedInformers().KubernetesClientSet
	c.cset = clientset.NewForConfigOrDie(arguments.GetRestConfig())
}

func (c *RefineryControllerImpl) Get(namespace, name string) (*v1alpha1.Refinery, error) {
	return c.lister.Refineries(namespace).Get(name)
}
