/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package fake

import (
	v1alpha1 "github.com/lwolf/kube-replay/pkg/apis/replay/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSilos implements SiloInterface
type FakeSilos struct {
	Fake *FakeKubereplayV1alpha1
	ns   string
}

var silosResource = schema.GroupVersionResource{Group: "kubereplay.lwolf.org", Version: "v1alpha1", Resource: "silos"}

var silosKind = schema.GroupVersionKind{Group: "kubereplay.lwolf.org", Version: "v1alpha1", Kind: "Silo"}

// Get takes name of the silo, and returns the corresponding silo object, and an error if there is any.
func (c *FakeSilos) Get(name string, options v1.GetOptions) (result *v1alpha1.Silo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(silosResource, c.ns, name), &v1alpha1.Silo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Silo), err
}

// List takes label and field selectors, and returns the list of Silos that match those selectors.
func (c *FakeSilos) List(opts v1.ListOptions) (result *v1alpha1.SiloList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(silosResource, silosKind, c.ns, opts), &v1alpha1.SiloList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SiloList{}
	for _, item := range obj.(*v1alpha1.SiloList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested silos.
func (c *FakeSilos) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(silosResource, c.ns, opts))

}

// Create takes the representation of a silo and creates it.  Returns the server's representation of the silo, and an error, if there is any.
func (c *FakeSilos) Create(silo *v1alpha1.Silo) (result *v1alpha1.Silo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(silosResource, c.ns, silo), &v1alpha1.Silo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Silo), err
}

// Update takes the representation of a silo and updates it. Returns the server's representation of the silo, and an error, if there is any.
func (c *FakeSilos) Update(silo *v1alpha1.Silo) (result *v1alpha1.Silo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(silosResource, c.ns, silo), &v1alpha1.Silo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Silo), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSilos) UpdateStatus(silo *v1alpha1.Silo) (*v1alpha1.Silo, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(silosResource, "status", c.ns, silo), &v1alpha1.Silo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Silo), err
}

// Delete takes name of the silo and deletes it. Returns an error if one occurs.
func (c *FakeSilos) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(silosResource, c.ns, name), &v1alpha1.Silo{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSilos) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(silosResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.SiloList{})
	return err
}

// Patch applies the patch and returns the patched silo.
func (c *FakeSilos) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Silo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(silosResource, c.ns, name, data, subresources...), &v1alpha1.Silo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Silo), err
}
