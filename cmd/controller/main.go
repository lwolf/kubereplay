package main

import (
	"flag"
	"k8s.io/client-go/util/workqueue"
	"fmt"
	"log"
	"time"
	"reflect"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/api/core/v1"
	"github.com/lwolf/kube-replay/pkg/signals"
	"github.com/lwolf/kube-replay/pkg/apis/replay/v1alpha1"
	"k8s.io/client-go/util/retry"
	client "github.com/lwolf/kube-replay/pkg/client/clientset/versioned"
	factory "github.com/lwolf/kube-replay/pkg/client/informers/externalversions"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/apimachinery/pkg/util/runtime"
	"github.com/lwolf/kube-replay/cmd/controller/utils"
	"github.com/lwolf/kube-replay/pkg/apis/replay"
)

var (

	apiserverURL = flag.String("apiserver", "", "Optional URL used to access the Kubernetes API server")
	// kubeconfig is the URL of the API server to connect to
	kubeconfig = flag.String("kubeconfig", "", "Optional kubeconfig path used to access the Kubernetes API server")

	// queue is a queue of resources to be processed. It performs exponential
	// backoff rate limiting, with a minimum retry period of 5 seconds and a
	// maximum of 1 minute.
	queue = workqueue.NewRateLimitingQueue(workqueue.NewItemExponentialFailureRateLimiter(time.Second*5, time.Minute))

	// stopCh can be used to stop all the informer, as well as control loops
	// within the application.
	stopCh = make(chan struct{})

	// sharedFactory is a shared informer factory that is used a a cache for
	// items in the API server. It saves each informer listing and watching the
	// same resources independently of each other, thus providing more up to
	// date results with less 'effort'
	sharedFactory factory.SharedInformerFactory

	// cl is a Kubernetes API client for our custom resource definition type
	cl client.Interface

	// kc is a Kubernetes API client for default resources
	kc k8sclient.Interface

)

func main() {

	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(*apiserverURL, *kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

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

	// we use a shared informer from the informer factory, to save calls to the
	// API as we grow our application and so state is consistent between our
	// control loops. We set a resync period of 30 seconds, in case any
	// create/replace/update/delete operations are missed when watching
	sharedFactory = factory.NewSharedInformerFactory(cl, time.Second*30)

	informer := sharedFactory.Kubereplay().V1alpha1().Silos().Informer()
	// we add a new event handler, watching for changes to API resources.
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: enqueue,
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					enqueue(cur)
				}
			},
			DeleteFunc: enqueue,
		},
	)
	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()
	// start the informer. This will cause it to begin receiving updates from
	// the configured API server and firing event handlers in response.
	sharedFactory.Start(stopCh)
	log.Printf("Started informer factory.")

	// wait for the informe rcache to finish performing it's initial sync of
	// resources
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		log.Fatalf("error waiting for informer cache to sync: %s", err.Error())
	}

	log.Printf("Finished populating shared informer cache.")
	// here we start just one worker reading objects off the queue. If you
	// wanted to parallelize this, you could start many instances of the worker
	// function, then ensure your application handles concurrency correctly.
	work()

}

// sync will attempt to 'Sync' an alert resource. It checks to see if the alert
// has already been sent, and if not will send it and update the resource
// accordingly. This method is called whenever this controller starts, and
// whenever the resource changes, and also periodically every resyncPeriod.
func sync(silo *v1alpha1.Silo) error {
	log.Printf("Found new event about silo '%s/%s'", silo.Namespace, silo.Name)
	// deploy new instance of goreplay for each silo without deployed status
	if silo.Status.Deployed != true {
		deploymentsClient := kc.AppsV1().Deployments(apiv1.NamespaceDefault)
		var spec replay.SiloSpec
		err := v1alpha1.Convert_v1alpha1_SiloSpec_To_replay_SiloSpec(&silo.Spec, &spec, nil)
		if err != nil {
			log.Fatalf("Unable to convert silo spec v1")
		}
		deployment := utils.CreateDeployment(silo.Name, &spec)

		// Create Deployment
		log.Printf("Creating silo deployment...")
		result, err := deploymentsClient.Create(deployment)
		if err != nil {
			panic(err)
		}
		log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Retrieve the latest version of Deployment before attempting update
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver

			silosClient := cl.KubereplayV1alpha1().Silos(silo.Namespace)

			result, getErr := silosClient.Get(silo.Name, metav1.GetOptions{})
			if getErr != nil {
				panic(fmt.Errorf("Failed to get latest version of Silo: %v", getErr))
			}
			result.Status.Deployed = true
			_, updateErr := silosClient.Update(result)
			return updateErr
		})
		if retryErr != nil {
			panic(fmt.Errorf("Update failed: %v", retryErr))
		}
		log.Printf("Updated deployment...")
	}
	return nil
}


func work() {
	log.Println("Starting processing the queue")
	for {
		// we read a message off the queue
		key, shutdown := queue.Get()

		// if the queue has been shut down, we should exit the work queue here
		if shutdown {
			stopCh <- struct{}{}
			return
		}

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

			// retrieve the latest version in the cache of this alert
			obj, err := sharedFactory.Kubereplay().V1alpha1().Silos().Lister().Silos(namespace).Get(name)

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
func enqueue(obj interface{}) {
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

func int32Ptr(i int32) *int32 { return &i }