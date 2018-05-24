package harvester

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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
	appsv1beta "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1lister "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
)

const controllerAgentName = "kubereplay-harvester-controller"

func (bc *HarvesterController) reconcileDeployment(green *appsv1beta.Deployment, blue *appsv1beta.Deployment, blueReplicas int32, greenReplicas int32) {
	log.Printf("reconciling deployment %s to %d/%d", green.Name, blueReplicas, greenReplicas)
	if *blue.Spec.Replicas != blueReplicas {
		log.Printf("blue replica needs reconcilation %d != %d", *blue.Spec.Replicas, blueReplicas)
		deploy, err := bc.kubernetesclient.ExtensionsV1beta1().Deployments(blue.Namespace).Get(blue.Name, metav1.GetOptions{})
		if err != nil {
			log.Printf("failed to get scale for deployment %s: %v", blue.Name, err)
		}
		deploy.Spec.Replicas = &blueReplicas
		deploy.Annotations[helpers.AnnotationKeyReplicas] = fmt.Sprintf("%d", blueReplicas)
		_, err = bc.kubernetesclient.ExtensionsV1beta1().Deployments(blue.Namespace).Update(deploy)
		if err != nil {
			log.Printf("failed to scale deployment %s to %d replicas: %v", blue.Name, blueReplicas, err)
		}
	}
	if *green.Spec.Replicas != greenReplicas {
		log.Printf("green replica needs reconcilation %d != %d", *green.Spec.Replicas, greenReplicas)
		deploy, err := bc.kubernetesclient.ExtensionsV1beta1().Deployments(green.Namespace).Get(green.Name, metav1.GetOptions{})
		if err != nil {
			log.Printf("failed to get scale for deployment %s: %v", green.Name, err)
		}
		deploy.Spec.Replicas = &greenReplicas
		deploy.Annotations[helpers.AnnotationKeyReplicas] = fmt.Sprintf("%d", greenReplicas)
		_, err = bc.kubernetesclient.ExtensionsV1beta1().Deployments(green.Namespace).Update(deploy)
		if err != nil {
			log.Printf("failed to scale deployment %s to %d replicas: %v", green.Name, greenReplicas, err)
		}
	}

}

func arrayToMap(deployments []*appsv1beta.Deployment) map[string]*appsv1beta.Deployment {
	res := make(map[string]*appsv1beta.Deployment)
	for _, d := range deployments {
		res[d.Name] = d
	}
	return res
}

func findOrphans(deployments map[string]*appsv1beta.Deployment) []string {
	var orphans []string
	for _, d := range deployments {
		if strings.HasSuffix(d.Name, "-gor") {
			_, ok := deployments[strings.TrimSuffix(d.Name, "-gor")]
			if !ok {
				orphans = append(orphans, d.Name)
			}
		}
	}
	return orphans
}

func (bc *HarvesterController) Reconcile(k types.ReconcileKey) error {
	log.Printf("running reconcile Harvester for %s", k.Name)
	h, err := bc.Get(k.Namespace, k.Name)
	if err != nil {
		return err
	}

	selector, err := metav1.LabelSelectorAsSelector(
		&metav1.LabelSelector{MatchLabels: h.Spec.Selector},
	)
	deploys, err := bc.deploymentLister.List(selector)
	if err != nil {
		return err
	}
	deploysMap := arrayToMap(deploys)
	orphans := findOrphans(deploysMap)
	propagationPolicy := metav1.DeletePropagationBackground
	for _, orph := range orphans {
		err = bc.kubernetesclient.ExtensionsV1beta1().Deployments(k.Namespace).Delete(
			orph, &metav1.DeleteOptions{PropagationPolicy: &propagationPolicy})
		if err != nil {
			log.Printf("failed to delete shadow deployment %s: %v", orph, err)
		}

	}

	var forceReconcile bool
	if h.Spec.SegmentSize != h.Status.SegmentSize {
		forceReconcile = true
	}

	for _, d := range deploys {
		a, ok := d.Annotations[helpers.AnnotationKeyDefault]

		if !ok {
			// annotation is not present, skipping
			continue
		}
		if a == helpers.AnnotationValueCapture {
			continue
		}
		blueName, ok := d.Annotations[helpers.AnnotationKeyShadow]
		if !ok {
			log.Printf("deployment %s does not have a shadow", d.Name)
			continue
		}
		blue, err := bc.deploymentLister.Deployments(d.Namespace).Get(blueName)
		if err != nil {
			log.Printf("failed to get deployment by shadow name %s: %v", blueName, err)
			continue
		}
		var blueReplicas, greenReplicas int32
		if forceReconcile {
			blueReplicas, greenReplicas = helpers.BlueGreenReplicas(*d.Spec.Replicas+*blue.Spec.Replicas, int32(h.Spec.SegmentSize))
		} else {
			ar, ok := d.Annotations[helpers.AnnotationKeyReplicas]
			if ok {
				v, err := strconv.Atoi(ar)
				if err == nil {
					if *d.Spec.Replicas == int32(v) {
						continue
					}
				}
			}
			blueReplicas, greenReplicas = helpers.BlueGreenReplicas(*d.Spec.Replicas, int32(h.Spec.SegmentSize))
		}
		log.Printf("new replicas count %d, %d. is forced reconcile=%v", blueReplicas, greenReplicas, forceReconcile)
		go bc.reconcileDeployment(d, blue, blueReplicas, greenReplicas)
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Harvester before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver

		result, getErr := bc.Get(h.Namespace, h.Name)
		if getErr != nil {
			log.Fatalf("Failed to get latest version of Harvester: %v", getErr)
		}
		result.Status.SegmentSize = h.Spec.SegmentSize
		_, updateErr := bc.harvesterclient.Harvesters(h.Namespace).Update(h)
		return updateErr
	})
	if retryErr != nil {
		log.Printf("Update failed: %v", retryErr)
		return retryErr
	}

	log.Printf("Finished processing harvester...")
	return nil
}

func (bc *HarvesterController) Lookup(k types.ReconcileKey) (interface{}, error) {
	return bc.harvesterLister.Harvesters(k.Namespace).Get(k.Name)
}

func (bc *HarvesterController) Get(namespace, name string) (*kubereplayv1alpha1.Harvester, error) {
	return bc.harvesterLister.Harvesters(namespace).Get(name)
}

// +controller:group=kubereplay,version=v1alpha1,kind=Harvester,resource=harvesters
// +informers:group=apps,version=v1,kind=Deployment
// +rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
type HarvesterController struct {
	args.InjectArgs
	harvesterLister kubereplayv1alpha1lister.HarvesterLister
	harvesterclient kubereplayv1alpha1client.KubereplayV1alpha1Interface

	deploymentLister appsv1lister.DeploymentLister
	kubernetesclient *kubernetes.Clientset

	recorder record.EventRecorder
}

// ProvideController provides a controller that will be run at startup.  Kubebuilder will use codegeneration
// to automatically register this controller in the inject package
func ProvideController(arguments args.InjectArgs) (*controller.GenericController, error) {
	bc := &HarvesterController{
		InjectArgs:       arguments,
		harvesterLister:  arguments.ControllerManager.GetInformerProvider(&kubereplayv1alpha1.Harvester{}).(kubereplayv1alpha1informer.HarvesterInformer).Lister(),
		harvesterclient:  arguments.Clientset.KubereplayV1alpha1(),
		deploymentLister: arguments.KubernetesInformers.Apps().V1().Deployments().Lister(),
		kubernetesclient: arguments.KubernetesClientSet,
		recorder:         arguments.CreateRecorder(controllerAgentName),
	}

	// Create a new controller that will call HarvesterController.Reconcile on changes to Harvesters
	gc := &controller.GenericController{
		Name:             "HarvesterController",
		Reconcile:        bc.Reconcile,
		InformerRegistry: arguments.ControllerManager,
	}
	if err := gc.Watch(&kubereplayv1alpha1.Harvester{}); err != nil {
		return gc, err
	}

	// INSERT ADDITIONAL WATCHES HERE BY CALLING gc.Watch.*() FUNCTIONS
	// NOTE: Informers for Kubernetes resources *MUST* be registered in the pkg/inject package so that they are started.
	if err := gc.WatchControllerOf(&appsv1beta.Deployment{}, eventhandlers.Path{bc.Lookup},
		predicates.ResourceVersionChanged); err != nil {
		return gc, err
	}
	return gc, nil
}
