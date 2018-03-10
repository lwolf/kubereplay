package fake

import (
	"github.com/lwolf/kubereplay/pkg/client/clientset_generated/clientset/typed/kubereplay/v1alpha1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/testing"
)

type FakeKubereplayV1alpha1 struct {
	*testing.Fake
}

func (c *FakeKubereplayV1alpha1) Harvesters(namespace string) v1alpha1.HarvesterInterface {
	return &FakeHarvesters{c, namespace}
}

func (c *FakeKubereplayV1alpha1) Refineries(namespace string) v1alpha1.RefineryInterface {
	return &FakeRefineries{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeKubereplayV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
