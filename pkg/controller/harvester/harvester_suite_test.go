package harvester_test

import (
	"testing"

	"github.com/kubernetes-sigs/kubebuilder/pkg/test"
	"k8s.io/client-go/rest"

	"github.com/lwolf/kubereplay/pkg/apis"
	"github.com/lwolf/kubereplay/pkg/client/clientset_generated/clientset"
	"github.com/lwolf/kubereplay/pkg/controller/harvester"
	"github.com/lwolf/kubereplay/pkg/controller/sharedinformers"
)

var testenv *test.TestEnvironment
var config *rest.Config
var cs *clientset.Clientset
var shutdown chan struct{}
var controller *harvester.HarvesterController
var si *sharedinformers.SharedInformers

func TestHarvester(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Harvester Suite", []Reporter{test.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	testenv = &test.TestEnvironment{CRDs: apis.APIMeta.GetCRDs()}
	var err error
	config, err = testenv.Start()
	Expect(err).NotTo(HaveOccurred())
	cs = clientset.NewForConfigOrDie(config)

	shutdown = make(chan struct{})
	si = sharedinformers.NewSharedInformers(config, shutdown)
	controller = harvester.NewHarvesterController(config, si)
	controller.Run(shutdown)
})

var _ = AfterSuite(func() {
	testenv.Stop()
})
