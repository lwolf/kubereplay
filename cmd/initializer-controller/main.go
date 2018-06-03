package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/lwolf/kubereplay/helpers"
	"github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	"github.com/lwolf/kubereplay/pkg/client/clientset/versioned"
	"github.com/lwolf/kubereplay/pkg/client/informers/externalversions"
	kubereplayv1alpha1lister "github.com/lwolf/kubereplay/pkg/client/listers/kubereplay/v1alpha1"
	"k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
)

const (
	defaultInitializerName = "kubereplay.initializer.lwolf.org"
)

var (
	initializerName string
	err             error
	clusterConfig   *rest.Config
	kubeconfig      string
)

var harvesterGVK = v1alpha1.SchemeGroupVersion.WithKind("Harvester")

func GenerateSidecar(refinerySvc string, port uint32) *corev1.Container {
	return &corev1.Container{
		Name:  "goreplay",
		Image: "buger/goreplay:latest",
		Args: []string{
			"--input-raw",
			fmt.Sprintf(":%d", port),
			"--output-tcp",
			fmt.Sprintf("%s:28020", refinerySvc),
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("10m"),
				corev1.ResourceMemory: resource.MustParse("64Mi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("10m"),
				corev1.ResourceMemory: resource.MustParse("64Mi"),
			},
		},
	}
}

func createShadowDeployment(d *v1beta1.Deployment, clientset *kubernetes.Clientset) {
	_, err = clientset.AppsV1beta1().Deployments(d.Namespace).Create(d)
	if err != nil {
		log.Printf("failed to create blue deployment %s: %v", d.Name, err)
	}

}

func initializeDeployment(deployment *v1beta1.Deployment, clientset *kubernetes.Clientset, lister kubereplayv1alpha1lister.HarvesterLister) error {
	if deployment.ObjectMeta.GetInitializers() != nil {
		pendingInitializers := deployment.ObjectMeta.GetInitializers().Pending
		if initializerName == pendingInitializers[0].Name {
			log.Printf("Initializing deployment: %s", deployment.Name)

			initializedDeploymentGreen := deployment.DeepCopy()

			// Remove self from the list of pending Initializers while preserving ordering.
			if len(pendingInitializers) == 1 {
				initializedDeploymentGreen.ObjectMeta.Initializers = nil
			} else {
				initializedDeploymentGreen.ObjectMeta.Initializers.Pending = append(pendingInitializers[:0], pendingInitializers[1:]...)
			}

			selector, err := metav1.LabelSelectorAsSelector(
				&metav1.LabelSelector{MatchLabels: deployment.ObjectMeta.GetLabels()},
			)
			var skip bool
			harvesters, err := lister.Harvesters(deployment.Namespace).List(labels.Everything())
			if err != nil {
				log.Printf("failed to get list of harvesters: %v", err)
				skip = true
			}
			var harvester *v1alpha1.Harvester
			for _, h := range harvesters {
				if labels.Equals(h.Spec.Selector, deployment.ObjectMeta.GetLabels()) {
					harvester = h
				}
			}

			if harvester == nil {
				log.Printf("debug: harvesters not found for deployment %s with selectors %v", deployment.Name, selector)
				skip = true
			}
			annotations := deployment.GetAnnotations()
			_, ok := annotations[helpers.AnnotationKeyDefault]
			if ok {
				skip = true
			}

			if skip {
				// Releasing original deployment
				_, err := clientset.AppsV1beta1().Deployments(deployment.Namespace).Update(initializedDeploymentGreen)
				if err != nil {
					log.Printf("failed to update initialized green deployment %s: %v ", initializedDeploymentGreen.Name, err)
					return err
				}
				return nil
			}

			initializedDeploymentBlue := &v1beta1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						*metav1.NewControllerRef(harvester, harvesterGVK),
					},
					Name:      fmt.Sprintf("%s-gor", deployment.Name),
					Namespace: deployment.Namespace,
					Labels:    deployment.ObjectMeta.Labels,
				},
				Spec: *deployment.Spec.DeepCopy(),
			}

			//Remove self from the list of pending Initializers while preserving ordering.
			if len(pendingInitializers) == 1 {
				initializedDeploymentBlue.ObjectMeta.Initializers = nil
			} else {
				initializedDeploymentBlue.ObjectMeta.Initializers.Pending = append(pendingInitializers[:0], pendingInitializers[1:]...)
			}

			greenAnnotations := initializedDeploymentGreen.GetAnnotations()
			blueReplicas, greenReplicas := helpers.BlueGreenReplicas(*deployment.Spec.Replicas, int32(harvester.Spec.SegmentSize))
			if greenAnnotations == nil {
				greenAnnotations = make(map[string]string)
			}
			// set annotation for original deployment
			greenAnnotations[helpers.AnnotationKeyDefault] = helpers.AnnotationValueSkip
			greenAnnotations[helpers.AnnotationKeyReplicas] = strconv.Itoa(int(greenReplicas))
			greenAnnotations[helpers.AnnotationKeyShadow] = initializedDeploymentBlue.Name
			initializedDeploymentGreen.Annotations = greenAnnotations
			initializedDeploymentGreen.Spec.Replicas = &greenReplicas

			blueAnnotations := initializedDeploymentBlue.GetAnnotations()
			if blueAnnotations == nil {
				blueAnnotations = make(map[string]string)
			}
			blueAnnotations[helpers.AnnotationKeyDefault] = helpers.AnnotationValueCapture
			blueAnnotations[helpers.AnnotationKeyReplicas] = strconv.Itoa(int(blueReplicas))
			blueAnnotations[helpers.AnnotationKeyShadow] = initializedDeploymentGreen.Name
			initializedDeploymentBlue.Annotations = blueAnnotations
			initializedDeploymentBlue.Spec.Replicas = &blueReplicas
			initializedDeploymentBlue.Status = v1beta1.DeploymentStatus{}

			sidecar := GenerateSidecar(
				fmt.Sprintf("refinery-%s.default", harvester.Spec.Refinery),
				// todo: remove port from harvester spec, get it directly from deployment
				harvester.Spec.AppPort,
			)

			_, err = clientset.AppsV1beta1().Deployments(deployment.Namespace).Update(initializedDeploymentGreen)
			if err != nil {
				log.Printf("failed to simply update initialized and updated green deployment %s: %v ", initializedDeploymentGreen.Name, err)
				return err
			}

			// Modify the Deployment's Pod template to include the Gor container
			initializedDeploymentBlue.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, *sidecar)
			// Creating new deployment in a gorouting, otherwise it will block and timeout
			go createShadowDeployment(initializedDeploymentBlue, clientset)

		}
	}

	return nil
}

func main() {
	flag.StringVar(&initializerName, "initializer-name", defaultInitializerName, "The initializer name")
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig")
	flag.Parse()
	log.Println("Starting the Kubernetes initializer...")
	log.Printf("Initializer name set to: %s", initializerName)

	if kubeconfig != "" {
		kubeConfigLocation := filepath.Join(kubeconfig)
		clusterConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigLocation)
	} else {
		clusterConfig, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Watch uninitialized Deployments in all namespaces.
	restClient := clientset.AppsV1beta1().RESTClient()
	watchlist := cache.NewListWatchFromClient(restClient, "deployments", corev1.NamespaceAll, fields.Everything())

	stop := make(chan struct{})
	cl := versioned.NewForConfigOrDie(clusterConfig)
	si := externalversions.NewSharedInformerFactory(cl, 30*time.Second)
	go si.Kubereplay().V1alpha1().Harvesters().Informer().Run(stop)
	si.WaitForCacheSync(stop)

	lister := si.Kubereplay().V1alpha1().Harvesters().Lister()

	// Wrap the returned watchlist to workaround the inability to include
	// the `IncludeUninitialized` list option when setting up watch clients.
	includeUninitializedWatchlist := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.IncludeUninitialized = true
			return watchlist.List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.IncludeUninitialized = true
			return watchlist.Watch(options)
		},
	}

	resyncPeriod := 30 * time.Second

	_, controller := cache.NewInformer(includeUninitializedWatchlist, &v1beta1.Deployment{}, resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				err := initializeDeployment(obj.(*v1beta1.Deployment), clientset, lister)
				if err != nil {
					log.Println(err)
				}
			},
		},
	)
	go controller.Run(stop)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")
	close(stop)
}
