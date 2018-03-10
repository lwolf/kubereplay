package main

import (
	"flag"
	"log"

	controllerlib "github.com/kubernetes-sigs/kubebuilder/pkg/controller"
	"github.com/kubernetes-sigs/kubebuilder/pkg/install"
	extensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	"github.com/lwolf/kubereplay/pkg/apis"
	"github.com/lwolf/kubereplay/pkg/controller"
)

var kubeconfig = flag.String("kubeconfig", "", "path to kubeconfig")
var installCRDs = flag.Bool("install-crds", true, "install the CRDs used by the controller as part of startup")

// Controller-manager main.
func main() {
	flag.Parse()
	config, err := controllerlib.GetConfig(*kubeconfig)
	if err != nil {
		log.Fatalf("Could not create Config for talking to the apiserver: %v", err)
	}

	if *installCRDs {
		err = install.NewInstaller(config).Install(&InstallStrategy{crds: apis.APIMeta.GetCRDs()})
		if err != nil {
			log.Fatalf("Could not create CRDs: %v", err)
		}
	}

	// Start the controllers
	controllers, _ := controller.GetAllControllers(config)
	controllerlib.StartControllerManager(controllers...)

	// Blockforever
	select {}
}

type InstallStrategy struct {
	install.EmptyInstallStrategy
	crds []extensionsv1beta1.CustomResourceDefinition
}

func (s *InstallStrategy) GetCRDs() []extensionsv1beta1.CustomResourceDefinition {
	return s.crds
}
