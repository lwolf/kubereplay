package refinery

import (
	"log"

	"github.com/kubernetes-sigs/kubebuilder/pkg/controller"
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/types"

	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/eventhandlers"
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/predicates"
	"github.com/lwolf/kubereplay/helpers"
	kubereplayv1alpha1 "github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	kubereplayv1alpha1client "github.com/lwolf/kubereplay/pkg/client/clientset/versioned/typed/kubereplay/v1alpha1"
	kubereplayv1alpha1informer "github.com/lwolf/kubereplay/pkg/client/informers/externalversions/kubereplay/v1alpha1"
	kubereplayv1alpha1lister "github.com/lwolf/kubereplay/pkg/client/listers/kubereplay/v1alpha1"
	"github.com/lwolf/kubereplay/pkg/inject/args"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1lister "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
)

const controllerAgentName = "kubereplay-refinery-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	MessageResourceSynced = "Refinery synced successfully"
)

func (bc *RefineryController) Reconcile(k types.ReconcileKey) error {
	log.Printf("Running reconcile Refinery for %s\n", k.Name)

	r, err := bc.Get(k.Namespace, k.Name)
	if err != nil {
		return err
	}

	sClient := bc.kubernetesclient.CoreV1().Services(r.Namespace)
	service := helpers.GenerateService(r.Name, &r.Spec)
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

	dClient := bc.kubernetesclient.AppsV1().Deployments(k.Namespace)
	deployment := helpers.GenerateDeployment(k.Name, r)
	_, err = dClient.Get(deployment.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
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

		result, getErr := bc.Get(r.Namespace, r.Name)
		if getErr != nil {
			log.Fatalf("Failed to get latest version of Silo: %v", getErr)
		}
		result.Status.Deployed = true
		_, updateErr := bc.refineryclient.Refineries(result.Namespace).Update(result)
		return updateErr
	})
	if retryErr != nil {
		log.Printf("Update failed: %v", retryErr)
		return retryErr
	}

	bc.recorder.Event(r, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (bc *RefineryController) Get(namespace, name string) (*kubereplayv1alpha1.Refinery, error) {
	return bc.refineryLister.Refineries(namespace).Get(name)
}

func (bc *RefineryController) Lookup(k types.ReconcileKey) (interface{}, error) {
	return bc.refineryLister.Refineries(k.Namespace).Get(k.Name)
}

// +controller:group=kubereplay,version=v1alpha1,kind=Refinery,resource=refineries
// +informers:group=apps,version=v1,kind=Deployment
// +rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
type RefineryController struct {
	args.InjectArgs

	refineryLister   kubereplayv1alpha1lister.RefineryLister
	refineryclient   kubereplayv1alpha1client.KubereplayV1alpha1Interface
	deploymentLister appsv1lister.DeploymentLister
	kubernetesclient *kubernetes.Clientset

	recorder record.EventRecorder
}

// ProvideController provides a controller that will be run at startup.  Kubebuilder will use codegeneration
// to automatically register this controller in the inject package
func ProvideController(arguments args.InjectArgs) (*controller.GenericController, error) {
	bc := &RefineryController{
		InjectArgs:       arguments,
		refineryLister:   arguments.ControllerManager.GetInformerProvider(&kubereplayv1alpha1.Refinery{}).(kubereplayv1alpha1informer.RefineryInformer).Lister(),
		refineryclient:   arguments.Clientset.KubereplayV1alpha1(),
		deploymentLister: arguments.KubernetesInformers.Apps().V1().Deployments().Lister(),
		kubernetesclient: arguments.KubernetesClientSet,
		recorder:         arguments.CreateRecorder(controllerAgentName),
	}

	// Create a new controller that will call RefineryController.Reconcile on changes to Refinerys
	gc := &controller.GenericController{
		Name:             "RefineryController",
		Reconcile:        bc.Reconcile,
		InformerRegistry: arguments.ControllerManager,
	}
	if err := gc.Watch(&kubereplayv1alpha1.Refinery{}); err != nil {
		return gc, err
	}
	// INSERT ADDITIONAL WATCHES HERE BY CALLING gc.Watch.*() FUNCTIONS
	// NOTE: Informers for Kubernetes resources *MUST* be registered in the pkg/inject package so that they are started.
	if err := gc.WatchControllerOf(&appsv1.Deployment{}, eventhandlers.Path{bc.Lookup},
		predicates.ResourceVersionChanged); err != nil {
		return gc, err
	}

	return gc, nil
}
