package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"
	"github.com/ghodss/yaml"
	"github.com/mohae/deepcopy"
	"k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
)

const (
	defaultAnnotation      = "kubereplay.lwolf.org/mode"
	harvesterAnnotation    = "kubereplay.lwolf.org/harvester"
	defaultInitializerName = "kubereplay.initializer.lwolf.org"
	annotationValueCapture = "blue"
	annotationValueSkip    = "green"
)

var (
	initializerName string
	annotation      string
	external        bool
	err             error
	clusterConfig   *rest.Config
	kubeconfig      string
)

type config struct {
	Containers []corev1.Container
}

func configmapToConfig(configmap *corev1.ConfigMap) (*config, error) {
	var c config
	err := yaml.Unmarshal([]byte(configmap.Data["config"]), &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func initializeDeployment(deployment *v1beta1.Deployment, clientset *kubernetes.Clientset) error {
	log.Printf("%s: getInitializers for %s %v", time.Now(), deployment.Name, deployment.ObjectMeta.GetInitializers())
	if deployment.ObjectMeta.GetInitializers() != nil {
		pendingInitializers := deployment.ObjectMeta.GetInitializers().Pending
		if initializerName == pendingInitializers[0].Name {
			log.Printf("Initializing deployment: %s", deployment.Name)
			o := deepcopy.Copy(deployment)
			initializedDeployment := o.(*v1beta1.Deployment)

			// Remove self from the list of pending Initializers while preserving ordering.
			if len(pendingInitializers) == 1 {
				initializedDeployment.ObjectMeta.Initializers = nil
			} else {
				initializedDeployment.ObjectMeta.Initializers.Pending = append(pendingInitializers[:0], pendingInitializers[1:]...)
			}

			// Check required annotation
			annotations := deployment.ObjectMeta.GetAnnotations()
			a, ok := annotations[annotation]
			if !ok || a == annotationValueSkip {
				log.Printf("Required '%s' annotation missing or sidecar is not required. skipping container injection", annotation)
				_, err := clientset.AppsV1beta1().Deployments(deployment.Namespace).Update(initializedDeployment)
				if err != nil {
					return err
				}
				return nil
			}

			harvesterName, ok := annotations[harvesterAnnotation]
			if !ok {
				log.Printf("harvester annotation does not exist, skipping...")
				return nil
			}

			configmapName := fmt.Sprintf("%s-sidecar", harvesterName)

			cm, err := clientset.CoreV1().ConfigMaps(deployment.Namespace).Get(configmapName, metav1.GetOptions{})
			if err != nil {
				log.Fatal(err)
			}

			c, err := configmapToConfig(cm)
			if err != nil {
				log.Fatal(err)
			}

			// Modify the Deployment's Pod template to include the Envoy container
			// and configuration volume. Then patch the original deployment.
			initializedDeployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, c.Containers...)

			oldData, err := json.Marshal(deployment)
			if err != nil {
				return err
			}

			newData, err := json.Marshal(initializedDeployment)
			if err != nil {
				return err
			}

			patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1beta1.Deployment{})
			if err != nil {
				return err
			}

			_, err = clientset.AppsV1beta1().Deployments(deployment.Namespace).Patch(deployment.Name, types.StrategicMergePatchType, patchBytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	flag.StringVar(&annotation, "annotation", defaultAnnotation, "The annotation to trigger initialization")
	flag.StringVar(&initializerName, "initializer-name", defaultInitializerName, "The initializer name")
	flag.BoolVar(&external, "external", false, "Run initializer using configmap")
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to kubeconfig")
	flag.Parse()
	log.Println("Starting the Kubernetes initializer...")
	log.Printf("Initializer name set to: %s", initializerName)

	if external {
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
				err := initializeDeployment(obj.(*v1beta1.Deployment), clientset)
				if err != nil {
					log.Println(err)
				}
			},
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")
	close(stop)
}
