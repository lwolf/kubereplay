package v1alpha1_test

import (
	"testing"

	"github.com/kubernetes-sigs/kubebuilder/pkg/test"
	"k8s.io/client-go/rest"

	"github.com/lwolf/kubereplay/pkg/client/clientset/versioned"
	"github.com/lwolf/kubereplay/pkg/inject"
)

var testenv *test.TestEnvironment
var config *rest.Config
var cs *versioned.Clientset

func TestV1alpha1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "v1 Suite", []Reporter{test.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	testenv = &test.TestEnvironment{CRDs: inject.Injector.CRDs}

	var err error
	config, err = testenv.Start()
	Expect(err).NotTo(HaveOccurred())

	cs = versioned.NewForConfigOrDie(config)
})

var _ = AfterSuite(func() {
	testenv.Stop()
})
