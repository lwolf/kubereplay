package harvester

import (
	"log"

	"github.com/kubernetes-sigs/kubebuilder/pkg/builders"

	"github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	listers "github.com/lwolf/kubereplay/pkg/client/listers_generated/kubereplay/v1alpha1"
	"github.com/lwolf/kubereplay/pkg/controller/sharedinformers"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement controller logic for the Harvester resource API

// Reconcile handles enqueued messages
func (c *HarvesterControllerImpl) Reconcile(u *v1alpha1.Harvester) error {
	// INSERT YOUR CODE HERE - implement controller logic to reconcile observed and desired state of the object
	log.Printf("Running reconcile Harvester for %s\n", u.Name)
	return nil
}

// +controller:group=kubereplay,version=v1alpha1,kind=Harvester,resource=harvesters
type HarvesterControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about Harvester
	lister listers.HarvesterLister
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *HarvesterControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// INSERT YOUR CODE HERE - add logic for initializing the controller as needed

	// Use the lister for indexing harvesters labels
	c.lister = arguments.GetSharedInformers().Factory.Kubereplay().V1alpha1().Harvesters().Lister()

	// To watch other resource types, uncomment this function and replace Foo with the resource name to watch.
	// Must define the func FooToHarvester(i interface{}) (string, error) {} that returns the Harvester
	// "namespace/name"" to reconcile in response to the updated Foo
	// Note: To watch Kubernetes resources, you must also update the StartAdditionalInformers function in
	// pkg/controllers/sharedinformers/informers.go
	//
	// arguments.Watch("HarvesterFoo",
	//     arguments.GetSharedInformers().Factory.Bar().V1beta1().Bars().Informer(),
	//     c.FooToHarvester)
}

func (c *HarvesterControllerImpl) Get(namespace, name string) (*v1alpha1.Harvester, error) {
	return c.lister.Harvesters(namespace).Get(name)
}
