package harvester

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	k8sclient "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/util/workqueue"

	"github.com/lwolf/kube-replay/pkg/apis/replay/v1alpha1"
	client "github.com/lwolf/kube-replay/pkg/client/clientset/versioned"
	factory "github.com/lwolf/kube-replay/pkg/client/informers/externalversions"
)

var (
	// queue is a queue of resources to be processed. It performs exponential
	// backoff rate limiting, with a minimum retry period of 5 seconds and a
	// maximum of 1 minute.
	queue = workqueue.NewRateLimitingQueue(workqueue.NewItemExponentialFailureRateLimiter(time.Second*5, time.Minute))
	// cl is a Kubernetes API client for our custom resource definition type
	cl client.Interface

	// kc is a Kubernetes API client for default resources
	kc k8sclient.Interface

	err error
)

func labelSelector(selector *map[string]string) string {
	// "app=kubereplay,module=test"
	var result []string
	for key, value := range *selector {
		l := fmt.Sprintf("%s=%s", key, value)
		result = append(result, l)
	}
	return strings.Join(result, ",")
}

func sync(r *v1alpha1.Harvester) error {
	log.Printf("Found new event about Harvester '%s/%s'", r.Namespace, r.Name)

	deployClient := kc.AppsV1().Deployments(apiv1.NamespaceDefault)

	selector := labelSelector(&r.Spec.Selector)
	if selector == "" {
		log.Printf("Empty selector found in %s Harvester", r.Name)
		return errors.New("empty selector found in Harvester")
	}

	deploys, err := deployClient.List(metav1.ListOptions{
		LabelSelector: selector,
	})
	log.Println("Trying to get list of deploys")
	if err != nil {
		log.Printf("Failed to get list of deploys with labels")
	}

	var rsNamesToProcess []string

	for _, d := range deploys.Items {
		log.Printf("Found RC %s", d.Name)
		alreadyProcessed := false
		for _, c := range d.Spec.Template.Spec.Containers {
			if c.Name == "goreplay" {
				alreadyProcessed = true
			}
		}
		if alreadyProcessed != true {
			rsNamesToProcess = append(rsNamesToProcess, d.Name)
		}

	}

	for _, name := range rsNamesToProcess {
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Retrieve the latest version of Deployment before attempting update
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver

			rc, getErr := deployClient.Get(name, metav1.GetOptions{})
			if getErr != nil {
				log.Fatalf("Failed to get latest version of RC: %v", getErr)
			}
			args := []string{
				"--input-raw",
				fmt.Sprintf(":%d", r.Spec.AppPort),
				"--output-tcp",
				fmt.Sprintf("refinery-%s.%s:28020", r.Spec.Refinery, r.Namespace),
			}

			gorContainer := apiv1.Container{
				Name:  "goreplay",
				Image: "buger/goreplay:latest",
				Args:  args,
			}

			rc.Spec.Template.Spec.Containers = append(rc.Spec.Template.Spec.Containers, gorContainer)

			_, updateErr := deployClient.Update(rc)
			return updateErr
		})
		if retryErr != nil {
			log.Printf("Update failed: %v", retryErr)
			return retryErr
		}
		log.Printf("RC updated...")

	}
	return nil
}

func Work(sharedFactory factory.SharedInformerFactory, cfg *restclient.Config) {
	log.Println("Starting processing the queue")
	// create an instance of our own API client
	cl, err = client.NewForConfig(cfg)

	if err != nil {
		log.Fatalf("Error creating custom api client: %s", err.Error())
	}

	log.Printf("Custom Kubernetes client created.")

	kc, err = k8sclient.NewForConfig(cfg)

	if err != nil {
		log.Fatalf("Error creating k8s api client: %s", err.Error())
	}

	log.Printf("Original Kubernetes client created.")

	for {
		// we read a message off the queue
		key, _ := queue.Get()
		//key, shutdown := queue.Get()

		//// if the queue has been shut down, we should exit the work queue here
		//if shutdown {
		//	stopCh <- struct{}{}
		//	return
		//}

		// convert the queue item into a string. If it's not a string, we'll
		// simply discard it as invalid data and log a message.
		var strKey string
		var ok bool
		if strKey, ok = key.(string); !ok {
			runtime.HandleError(fmt.Errorf("key in queue should be of type string but got %T. discarding", key))
			return
		}

		// we define a function here to process a queue item, so that we can
		// use 'defer' to make sure the message is marked as Done on the queue
		func(key string) {
			defer queue.Done(key)

			// attempt to split the 'key' into namespace and object name
			namespace, name, err := cache.SplitMetaNamespaceKey(strKey)

			if err != nil {
				runtime.HandleError(fmt.Errorf("error splitting meta namespace key into parts: %s", err.Error()))
				return
			}

			log.Printf("Read item '%s/%s' off workqueue. Processing...", namespace, name)

			// retrieve the latest version in the cache of this refinery
			obj, err := sharedFactory.Kubereplay().V1alpha1().Harvesters().Lister().Harvesters(namespace).Get(name)
			if err != nil {
				runtime.HandleError(fmt.Errorf("error getting object '%s/%s' from api: %s", namespace, name, err.Error()))
				return
			}

			log.Printf("Got most up to date version of '%s/%s'. Syncing...", namespace, name)

			// attempt to sync the current state of the world with the desired!
			// If sync returns an error, we skip calling `queue.Forget`,
			// thus causing the resource to be requeued at a later time.
			if err := sync(obj); err != nil {
				runtime.HandleError(fmt.Errorf("error processing item '%s/%s': %s", namespace, name, err.Error()))
				return
			}

			log.Printf("Finished processing '%s/%s' successfully! Removing from queue.", namespace, name)

			// as we managed to process this successfully, we can forget it
			// from the work queue altogether.
			queue.Forget(key)
		}(strKey)
	}
}

// enqueue will add an object 'obj' into the workqueue. The object being added
// must be of type metav1.Object, metav1.ObjectAccessor or cache.ExplicitKey.
func Enqueue(obj interface{}) {
	// DeletionHandlingMetaNamespaceKeyFunc will convert an object into a
	// 'namespace/name' string. We do this because our item may be processed
	// much later than now, and so we want to ensure it gets a fresh copy of
	// the resource when it starts. Also, this allows us to keep adding the
	// same item into the work queue without duplicates building up.
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("error obtaining key for object being enqueue: %s", err.Error()))
		return
	}
	// add the item to the queue
	queue.Add(key)
}
