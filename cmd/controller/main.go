package main

import (
	"flag"
	"log"
	"reflect"
	"time"

	"github.com/lwolf/kube-replay/cmd/controller/harvester"
	"github.com/lwolf/kube-replay/cmd/controller/refinery"
	client "github.com/lwolf/kube-replay/pkg/client/clientset/versioned"
	factory "github.com/lwolf/kube-replay/pkg/client/informers/externalversions"
	"github.com/lwolf/kube-replay/pkg/signals"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	apiserverURL = flag.String("apiserver", "", "Optional URL used to access the Kubernetes API server")
	// kubeconfig is the URL of the API server to connect to
	kubeconfig = flag.String("kubeconfig", "", "Optional kubeconfig path used to access the Kubernetes API server")

	// sharedFactory is a shared informer factory that is used a a cache for
	// items in the API server. It saves each informer listing and watching the
	// same resources independently of each other, thus providing more up to
	// date results with less 'effort'
	sharedFactory factory.SharedInformerFactory

	// cl is a Kubernetes API client for our custom resource definition type
	cl client.Interface
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

	// we use a shared informer from the informer factory, to save calls to the
	// API as we grow our application and so state is consistent between our
	// control loops. We set a resync period of 30 seconds, in case any
	// create/replace/update/delete operations are missed when watching
	sharedFactory = factory.NewSharedInformerFactory(cl, time.Second*30)

	informer := sharedFactory.Kubereplay().V1alpha1().Refineries().Informer()
	// we add a new event handler, watching for changes to API resources.
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: refinery.Enqueue,
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					refinery.Enqueue(cur)
				}
			},
			DeleteFunc: refinery.Enqueue,
		},
	)

	hInformer := sharedFactory.Kubereplay().V1alpha1().Harvesters().Informer()

	hInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: harvester.Enqueue,
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					harvester.Enqueue(cur)
				}
			},
			DeleteFunc: harvester.Enqueue,
		},
	)

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()
	// start the informer. This will cause it to begin receiving updates from
	// the configured API server and firing event handlers in response.
	sharedFactory.Start(stopCh)
	log.Printf("Started informer factory.")

	// wait for the informer cache to finish performing it's initial sync of
	// resources
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		log.Fatalf("error waiting for refinery informer cache to sync: %s", err.Error())
	}

	// wait for the informer cache to finish performing it's initial sync of
	// resources
	if !cache.WaitForCacheSync(stopCh, hInformer.HasSynced) {
		log.Fatalf("error waiting for harvester informer cache to sync: %s", err.Error())
	}

	log.Printf("Finished populating shared informers cache.")
	// here we start just one worker reading objects off the queue. If you
	// wanted to parallelize this, you could start many instances of the worker
	// function, then ensure your application handles concurrency correctly.
	// XXX: temp hack to run both watchers
	go refinery.Work(sharedFactory, cfg)
	harvester.Work(sharedFactory, cfg)
}
