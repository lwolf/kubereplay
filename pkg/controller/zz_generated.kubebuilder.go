package controller

import (
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller"
	"github.com/lwolf/kubereplay/pkg/controller/harvester"
	"github.com/lwolf/kubereplay/pkg/controller/refinery"
	"github.com/lwolf/kubereplay/pkg/controller/sharedinformers"
	"k8s.io/client-go/rest"
)

func GetAllControllers(config *rest.Config) ([]controller.Controller, chan struct{}) {
	shutdown := make(chan struct{})
	si := sharedinformers.NewSharedInformers(config, shutdown)
	return []controller.Controller{
		harvester.NewHarvesterController(config, si),
		refinery.NewRefineryController(config, si),
	}, shutdown
}
