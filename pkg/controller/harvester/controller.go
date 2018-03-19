package harvester

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kubernetes-sigs/kubebuilder/pkg/builders"
	extBeta "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	extv1listers "k8s.io/client-go/listers/extensions/v1beta1"

	"github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	listers "github.com/lwolf/kubereplay/pkg/client/listers_generated/kubereplay/v1alpha1"
	"github.com/lwolf/kubereplay/pkg/constants"
	"github.com/lwolf/kubereplay/pkg/controller/refinery"
	"github.com/lwolf/kubereplay/pkg/controller/sharedinformers"
)

// Created by "kubebuilder create resource" for you to implement controller logic for the Harvester resource API

func int32Ptr(i int32) *int32 { return &i }

// Reconcile handles enqueued messages
func (c *HarvesterControllerImpl) Reconcile(u *v1alpha1.Harvester) error {
	log.Printf("***********************")
	log.Printf("running reconcile Harvester for %s\n", u.Name)

	selector, err := metav1.LabelSelectorAsSelector(
		&metav1.LabelSelector{MatchLabels: u.Spec.Selector},
	)

	configName := refinery.ConfigmapName(u.Name)
	_, err = c.cs.CoreV1().ConfigMaps("default").Get(configName, metav1.GetOptions{})
	// TODO: properly handle http errors here
	if err != nil {
		cfg := refinery.GenerateConfigmap(configName, &u.Spec)
		_, err = c.cs.CoreV1().ConfigMaps(u.Namespace).Create(cfg)
		if err != nil {
			log.Printf("failed to create configmap %v", err)
			return err
		}
	}

	deploys, err := c.extDeploymentLister.List(selector)

	if err != nil {
		log.Printf("Failed to get list of deploys with labels")
		return err
	}

	for _, d := range deploys {
		log.Printf("%s current scaling %d", d.Name, *d.Spec.Replicas)
		a, ok := d.Annotations[constants.AnnotationKeyDefault]

		if !ok {
			// annotation is not present
			ownerReferences := []metav1.OwnerReference{
				{
					Name:       u.Name,
					UID:        u.UID,
					Kind:       "Harvester",
					APIVersion: v1alpha1.SchemeGroupVersion.String(),
				},
			}
			annotations := d.Annotations
			replicas := d.Spec.Replicas
			log.Printf("current annotations: %v", annotations)
			if annotations == nil {
				log.Printf("annotations is nil")
				annotations = map[string]string{}
			}
			// set annotation for original deployment
			annotations[constants.AnnotationKeyDefault] = constants.AnnotationValueSkip
			annotations[constants.AnnotationKeyHarvester] = u.Name
			d.ObjectMeta.OwnerReferences = ownerReferences
			d.Annotations = annotations
			d.Spec.Replicas = int32Ptr(0)

			//ownR := append(ownerReferences, metav1.OwnerReference{
			//	Name: d.Name,
			//	UID: d.UID,
			//	Kind: d.Kind,
			//	APIVersion: d.TypeMeta.APIVersion,
			//})
			//log.Println(d)
			//log.Println(ownR)

			blueDeploy := extBeta.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: ownerReferences,
					Name:            fmt.Sprintf("%s-gor", d.Name),
					Labels:          d.Labels,
					Annotations: map[string]string{
						constants.AnnotationKeyDefault:   constants.AnnotationValueCapture,
						constants.AnnotationKeyHarvester: u.Name,
						constants.AnnotationKeyMaster:    string(d.UID),
					},
				},
				Spec: *d.Spec.DeepCopy(),
			}
			blueDeploy.Spec.Replicas = replicas

			dc := c.cs.ExtensionsV1beta1().Deployments(u.Namespace)

			_, err = dc.Create(&blueDeploy)
			if err != nil {
				log.Printf("failed to create shadow deployment %v", err)
			}

			_, err = dc.Update(d)
			if err != nil {
				log.Printf("failed to update deployment %v", err)
				return err
			}
		} else {
			if a == constants.AnnotationValueCapture {
				log.Printf("resource %s is fully automated", d.Name)
				continue
			}
			// todo: handle scaling here
			log.Printf("debug: annotation is already present %s, need to check scale ", a)
			return nil
		}
		//if d.Annotations
		// 1. check that deploy is not processed yet
		// 2. clone deployment
		// 3. add annotations to both deployments
		// 4. update deployments
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

	cs *kubernetes.Clientset
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *HarvesterControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// INSERT YOUR CODE HERE - add logic for initializing the controller as needed

	// Use the lister for indexing harvesters labels
	c.lister = arguments.GetSharedInformers().Factory.Kubereplay().V1alpha1().Harvesters().Lister()
	c.extDeploymentLister = arguments.GetSharedInformers().KubernetesFactory.Extensions().V1beta1().Deployments().Lister()
	c.coreLister = arguments.GetSharedInformers().KubernetesFactory.Core().V1().ConfigMaps().Lister()
	c.cs = arguments.GetSharedInformers().KubernetesClientSet

	// To watch other resource types, uncomment this function and replace Foo with the resource name to watch.
	// Must define the func FooToHarvester(i interface{}) (string, error) {} that returns the Harvester
	// "namespace/name"" to reconcile in response to the updated Foo
	// Note: To watch Kubernetes resources, you must also update the StartAdditionalInformers function in
	// pkg/controllers/sharedinformers/informers.go
	//
	// arguments.Watch("HarvesterFoo",
	//     arguments.GetSharedInformers().Factory.Bar().V1beta1().Bars().Informer(),
	//     c.FooToHarvester)
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
	log.Printf("%s: ExtDeploymentToHarvesters %s", time.Now(), d.Name)

	harvesters, err := c.lister.List(labels.Everything())
	if err != nil {
		log.Printf("Failed to get list of harvesters: %v", err)
		return []string{}, nil
	}
	for _, h := range harvesters {
		if labels.Equals(d.Labels, h.Spec.Selector) {
			log.Printf("harvester with matching selector is found %s", h.Name)
			return []string{d.Namespace + "/" + h.Name}, nil
		}
	}

	log.Printf("%s: ExtDeploymentToHarvesters done processing ....", time.Now())
	return []string{}, nil
}
