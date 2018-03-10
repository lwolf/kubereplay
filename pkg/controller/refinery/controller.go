package refinery

import (
	"log"

	"github.com/kubernetes-sigs/kubebuilder/pkg/builders"
	v1listers "k8s.io/client-go/listers/apps/v1"

	"github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	"github.com/lwolf/kubereplay/pkg/client/clientset_generated/clientset"
	listers "github.com/lwolf/kubereplay/pkg/client/listers_generated/kubereplay/v1alpha1"
	"github.com/lwolf/kubereplay/pkg/controller/sharedinformers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement controller logic for the Refinery resource API

// Reconcile handles enqueued messages
func (c *RefineryControllerImpl) Reconcile(u *v1alpha1.Refinery) error {
	// INSERT YOUR CODE HERE - implement controller logic to reconcile observed and desired state of the object
	log.Printf("Running reconcile Refinery for %s\n", u.Name)
	if u.Status.Deployed == true {
		log.Printf("Refinery %s already processed, skipping\n", u.Name)
		return nil
	}

	deploymentsClient := c.cs.AppsV1().Deployments(u.Namespace)

	service := GenerateService(u.Name, &u.Spec)
	serviceClient := c.cs.CoreV1().Services(u.Namespace)
	_, err := serviceClient.Create(service)
	if err != nil {
		log.Printf("Failed to create service: %v", err)
		return err
	}

	deployment := GenerateDeployment(u.Name, u)
	// Create Deployment
	log.Printf("Creating refinery deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		log.Printf("Failed to create deployment: %v", err)
		return err
	}
	log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
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
	log.Printf("Deployment updated...")

	return nil
}

// +controller:group=kubereplay,version=v1alpha1,kind=Refinery,resource=refineries
type RefineryControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about Refinery
	lister listers.RefineryLister

	deploymentLister v1listers.DeploymentLister

	cset *clientset.Clientset
	cs   *kubernetes.Clientset
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *RefineryControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// INSERT YOUR CODE HERE - add logic for initializing the controller as needed

	// Use the lister for indexing refineries labels
	c.lister = arguments.GetSharedInformers().Factory.Kubereplay().V1alpha1().Refineries().Lister()
	c.deploymentLister = arguments.GetSharedInformers().KubernetesFactory.Apps().V1().Deployments().Lister()

	c.cs = arguments.GetSharedInformers().KubernetesClientSet
	c.cset = clientset.NewForConfigOrDie(arguments.GetRestConfig())

	// To watch other resource types, uncomment this function and replace Foo with the resource name to watch.
	// Must define the func FooToRefinery(i interface{}) (string, error) {} that returns the Refinery
	// "namespace/name"" to reconcile in response to the updated Foo
	// Note: To watch Kubernetes resources, you must also update the StartAdditionalInformers function in
	// pkg/controllers/sharedinformers/informers.go
	//
	// arguments.Watch("RefineryFoo",
	//     arguments.GetSharedInformers().Factory.Bar().V1beta1().Bars().Informer(),
	//     c.FooToRefinery)
}

func (c *RefineryControllerImpl) Get(namespace, name string) (*v1alpha1.Refinery, error) {
	return c.lister.Refineries(namespace).Get(name)
}
