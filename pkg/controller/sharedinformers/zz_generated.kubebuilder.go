package sharedinformers

import (
	"time"

	"github.com/kubernetes-sigs/kubebuilder/pkg/controller"
	"github.com/lwolf/kubereplay/pkg/client/clientset_generated/clientset"
	"github.com/lwolf/kubereplay/pkg/client/informers_generated/externalversions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// SharedInformers wraps all informers used by controllers so that
// they are shared across controller implementations
type SharedInformers struct {
	controller.SharedInformersDefaults
	Factory externalversions.SharedInformerFactory
}

// newSharedInformers returns a set of started informers
func NewSharedInformers(config *rest.Config, shutdown <-chan struct{}) *SharedInformers {
	si := &SharedInformers{
		controller.SharedInformersDefaults{},
		externalversions.NewSharedInformerFactory(clientset.NewForConfigOrDie(config), 10*time.Minute),
	}
	if si.SetupKubernetesTypes() {
		si.InitKubernetesInformers(config)
	}
	si.Init()
	si.startInformers(shutdown)
	si.StartAdditionalInformers(shutdown)
	return si
}

// startInformers starts all of the informers
func (si *SharedInformers) startInformers(shutdown <-chan struct{}) {
	go si.Factory.Kubereplay().V1alpha1().Refineries().Informer().Run(shutdown)
	go si.Factory.Kubereplay().V1alpha1().Harvesters().Informer().Run(shutdown)
}

// ControllerInitArguments are arguments provided to the Init function for a new controller.
type ControllerInitArguments interface {
	// GetSharedInformers returns the SharedInformers that can be used to access
	// informers and listers for watching and indexing Kubernetes Resources
	GetSharedInformers() *SharedInformers

	// GetRestConfig returns the Config to create new client-go clients
	GetRestConfig() *rest.Config

	// Watch uses resourceInformer to watch a resource.  When create, update, or deletes
	// to the resource type are encountered, watch uses watchResourceToReconcileResourceKey
	// to lookup the key for the resource reconciled by the controller (maybe a different type
	// than the watched resource), and enqueue it to be reconciled.
	// watchName: name of the informer.  may appear in logs
	// resourceInformer: gotten from the SharedInformer.  controls which resource type is watched
	// getReconcileKeys: takes an instance of the watched resource and returns
	//                   a slice of keys for the reconciled resource type to enqueue.
	Watch(watchName string, resourceInformer cache.SharedIndexInformer,
		getReconcileKeys func(interface{}) ([]string, error))
}

type ControllerInitArgumentsImpl struct {
	Si *SharedInformers
	Rc *rest.Config
	Rk func(key string) error
}

func (c ControllerInitArgumentsImpl) GetSharedInformers() *SharedInformers {
	return c.Si
}

func (c ControllerInitArgumentsImpl) GetRestConfig() *rest.Config {
	return c.Rc
}

// Watch uses resourceInformer to watch a resource.  When create, update, or deletes
// to the resource type are encountered, watch uses watchResourceToReconcileResourceKey
// to lookup the key for the resource reconciled by the controller (maybe a different type
// than the watched resource), and enqueue it to be reconciled.
// watchName: name of the informer.  may appear in logs
// resourceInformer: gotten from the SharedInformer.  controls which resource type is watched
// getReconcileKey: takes an instance of the watched resource and returns
//                                      a key for the reconciled resource type to enqueue.
func (c ControllerInitArgumentsImpl) Watch(
	watchName string, resourceInformer cache.SharedIndexInformer,
	getReconcileKey func(interface{}) ([]string, error)) {
	c.Si.Watch(watchName, resourceInformer, getReconcileKey, c.Rk)
}

type Controller interface{}

// ControllerInit new controllers should implement this.  It is more flexible in
// allowing additional options to be passed in
type ControllerInit interface {
	Init(args ControllerInitArguments)
}
