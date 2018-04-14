/*
Copyright 2017 Sergey Nuzhdin.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/package fake

import (
	v1alpha1 "github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRefineries implements RefineryInterface
type FakeRefineries struct {
	Fake *FakeKubereplayV1alpha1
	ns   string
}

var refineriesResource = schema.GroupVersionResource{Group: "kubereplay.lwolf.org", Version: "v1alpha1", Resource: "refineries"}

var refineriesKind = schema.GroupVersionKind{Group: "kubereplay.lwolf.org", Version: "v1alpha1", Kind: "Refinery"}

// Get takes name of the refinery, and returns the corresponding refinery object, and an error if there is any.
func (c *FakeRefineries) Get(name string, options v1.GetOptions) (result *v1alpha1.Refinery, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(refineriesResource, c.ns, name), &v1alpha1.Refinery{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Refinery), err
}

// List takes label and field selectors, and returns the list of Refineries that match those selectors.
func (c *FakeRefineries) List(opts v1.ListOptions) (result *v1alpha1.RefineryList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(refineriesResource, refineriesKind, c.ns, opts), &v1alpha1.RefineryList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.RefineryList{}
	for _, item := range obj.(*v1alpha1.RefineryList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested refineries.
func (c *FakeRefineries) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(refineriesResource, c.ns, opts))

}

// Create takes the representation of a refinery and creates it.  Returns the server's representation of the refinery, and an error, if there is any.
func (c *FakeRefineries) Create(refinery *v1alpha1.Refinery) (result *v1alpha1.Refinery, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(refineriesResource, c.ns, refinery), &v1alpha1.Refinery{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Refinery), err
}

// Update takes the representation of a refinery and updates it. Returns the server's representation of the refinery, and an error, if there is any.
func (c *FakeRefineries) Update(refinery *v1alpha1.Refinery) (result *v1alpha1.Refinery, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(refineriesResource, c.ns, refinery), &v1alpha1.Refinery{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Refinery), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRefineries) UpdateStatus(refinery *v1alpha1.Refinery) (*v1alpha1.Refinery, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(refineriesResource, "status", c.ns, refinery), &v1alpha1.Refinery{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Refinery), err
}

// Delete takes name of the refinery and deletes it. Returns an error if one occurs.
func (c *FakeRefineries) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(refineriesResource, c.ns, name), &v1alpha1.Refinery{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRefineries) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(refineriesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.RefineryList{})
	return err
}

// Patch applies the patch and returns the patched refinery.
func (c *FakeRefineries) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Refinery, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(refineriesResource, c.ns, name, data, subresources...), &v1alpha1.Refinery{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Refinery), err
}
