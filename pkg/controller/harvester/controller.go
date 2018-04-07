package harvester

import (
	"fmt"
	"log"
	"strconv"
	"github.com/kubernetes-sigs/kubebuilder/pkg/builders"
	extBeta "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	extv1listers "k8s.io/client-go/listers/extensions/v1beta1"

	"github.com/lwolf/kubereplay/helpers"
	"github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	"github.com/lwolf/kubereplay/pkg/client/clientset_generated/clientset"
	listers "github.com/lwolf/kubereplay/pkg/client/listers_generated/kubereplay/v1alpha1"
	"github.com/lwolf/kubereplay/pkg/constants"
	"github.com/lwolf/kubereplay/pkg/controller/sharedinformers"
	"k8s.io/client-go/util/retry"
)

func (c *HarvesterControllerImpl) reconcileDeployment(green *extBeta.Deployment, blue *extBeta.Deployment, blueReplicas int32, greenReplicas int32){
	log.Printf("reconciling deployment %s to %d/%d", green.Name, blueReplicas, greenReplicas)
	if *blue.Spec.Replicas != blueReplicas {
		log.Printf("blue replica needs reconcilation %d != %d", *blue.Spec.Replicas, blueReplicas)
		deploy, err := c.cs.ExtensionsV1beta1().Deployments(blue.Namespace).Get(blue.Name, metav1.GetOptions{})
		if err != nil {
			log.Printf("failed to get scale for deployment %s: %v", blue.Name, err)
		}
		deploy.Spec.Replicas = &blueReplicas
		deploy.Annotations[constants.AnnotationKeyReplicas] = fmt.Sprintf("%d", blueReplicas)
		_, err = c.cs.ExtensionsV1beta1().Deployments(blue.Namespace).Update(deploy)
		if err != nil {
			log.Printf("failed to scale deployment %s to %d replicas: %v", blue.Name, blueReplicas, err)
		}
	}
	if *green.Spec.Replicas != greenReplicas {
		log.Printf("green replica needs reconcilation %d != %d", *green.Spec.Replicas, greenReplicas)
		deploy, err := c.cs.ExtensionsV1beta1().Deployments(green.Namespace).Get(green.Name, metav1.GetOptions{})
		if err != nil {
			log.Printf("failed to get scale for deployment %s: %v", green.Name, err)
		}
		deploy.Spec.Replicas = &greenReplicas
		deploy.Annotations[constants.AnnotationKeyReplicas] = fmt.Sprintf("%d", greenReplicas)
		_, err = c.cs.ExtensionsV1beta1().Deployments(green.Namespace).Update(deploy)
		if err != nil {
			log.Printf("failed to scale deployment %s to %d replicas: %v", green.Name, greenReplicas, err)
		}
	}
}

// Reconcile handles enqueued messages
func (c *HarvesterControllerImpl) Reconcile(u *v1alpha1.Harvester) error {
	log.Printf("running reconcile Harvester for %s", u.Name)

	selector, err := metav1.LabelSelectorAsSelector(
		&metav1.LabelSelector{MatchLabels: u.Spec.Selector},
	)
	deploys, err := c.extDeploymentLister.List(selector)
	if err != nil {
		return err
	}
	var forceReconcile bool
	if u.Spec.SegmentSize != u.Status.SegmentSize {
		forceReconcile = true
	}

	for _, d := range deploys {
		a, ok := d.Annotations[constants.AnnotationKeyDefault]

		if !ok {
			// annotation is not present, skipping
			continue
		}
		if a == constants.AnnotationValueCapture {
			continue
		}
		blueName, ok := d.Annotations[constants.AnnotationKeyShadow]
		if !ok {
			log.Printf("deployment %s does not have a shadow", d.Name)
			continue
		}
		blue, err := c.extDeploymentLister.Deployments(d.Namespace).Get(blueName)
		if err != nil {
			log.Printf("failed to get deployment by shadow name %s: %v", blueName, err)
			continue
		}
		var blueReplicas, greenReplicas int32
		if forceReconcile {
			blueReplicas, greenReplicas = helpers.BlueGreenReplicas(*d.Spec.Replicas+*blue.Spec.Replicas, int32(u.Spec.SegmentSize))
		} else {
			ar, ok := d.Annotations[constants.AnnotationKeyReplicas]
			if ok {
				v, err := strconv.Atoi(ar)
				if err == nil {
					if *d.Spec.Replicas == int32(v) {
						continue
					}
				}
			}
			blueReplicas, greenReplicas = helpers.BlueGreenReplicas(*d.Spec.Replicas, int32(u.Spec.SegmentSize))
		}
		log.Printf("new replicas count %d, %d", blueReplicas, greenReplicas)
		go c.reconcileDeployment(d, blue, blueReplicas, greenReplicas)
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Harvester before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver

		result, getErr := c.Get(u.Namespace, u.Name)
		if getErr != nil {
			log.Fatalf("Failed to get latest version of Harvester: %v", getErr)
		}
		result.Status.SegmentSize = u.Spec.SegmentSize
		_, updateErr := c.client.KubereplayV1alpha1().Harvesters(u.Namespace).Update(u)
		return updateErr
	})
	if retryErr != nil {
		log.Printf("Update failed: %v", retryErr)
		return retryErr
	}

	log.Printf("Finished processing harvester...")

	return nil
}

// +controller:group=kubereplay,version=v1alpha1,kind=Harvester,resource=harvesters
type HarvesterControllerImpl struct {
	builders.DefaultControllerFns

	lister              listers.HarvesterLister
	extDeploymentLister extv1listers.DeploymentLister
	coreLister          corev1listers.ConfigMapLister

	cs     *kubernetes.Clientset
	client *clientset.Clientset
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *HarvesterControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// Use the lister for indexing harvesters labels
	c.lister = arguments.GetSharedInformers().Factory.Kubereplay().V1alpha1().Harvesters().Lister()
	c.extDeploymentLister = arguments.GetSharedInformers().KubernetesFactory.Extensions().V1beta1().Deployments().Lister()
	c.coreLister = arguments.GetSharedInformers().KubernetesFactory.Core().V1().ConfigMaps().Lister()
	c.cs = arguments.GetSharedInformers().KubernetesClientSet
	c.client = clientset.NewForConfigOrDie(arguments.GetRestConfig())

	arguments.Watch(
		"HarvesterExtDeployment",
		arguments.GetSharedInformers().KubernetesFactory.Extensions().V1beta1().Deployments().Informer(),
		c.ExtDeploymentToHarvesters)
}

func (c *HarvesterControllerImpl) Get(namespace, name string) (*v1alpha1.Harvester, error) {
	return c.lister.Harvesters(namespace).Get(name)
}

func (c *HarvesterControllerImpl) ExtDeploymentToHarvesters(i interface{}) ([]string, error) {
	d, _ := i.(*extBeta.Deployment)
	harvesters, err := c.lister.List(labels.Everything())
	if err != nil {
		log.Printf("failed to get list of harvesters: %v", err)
		return []string{}, nil
	}
	for _, h := range harvesters {
		if labels.Equals(d.Labels, h.Spec.Selector) {
			a, ok := d.Annotations[constants.AnnotationKeyDefault]
			if ok && a == constants.AnnotationValueCapture {
				continue
			}
			return []string{d.Namespace + "/" + h.Name}, nil
		}
	}

	return []string{}, nil
}
