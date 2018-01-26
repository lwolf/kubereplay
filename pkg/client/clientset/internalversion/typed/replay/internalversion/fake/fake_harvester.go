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
	replay "github.com/lwolf/kube-replay/pkg/apis/replay"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeHarvesters implements HarvesterInterface
type FakeHarvesters struct {
	Fake *FakeKubereplay
	ns   string
}

var harvestersResource = schema.GroupVersionResource{Group: "kubereplay.lwolf.org", Version: "", Resource: "harvesters"}

var harvestersKind = schema.GroupVersionKind{Group: "kubereplay.lwolf.org", Version: "", Kind: "Harvester"}

// Get takes name of the harvester, and returns the corresponding harvester object, and an error if there is any.
func (c *FakeHarvesters) Get(name string, options v1.GetOptions) (result *replay.Harvester, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(harvestersResource, c.ns, name), &replay.Harvester{})

	if obj == nil {
		return nil, err
	}
	return obj.(*replay.Harvester), err
}

// List takes label and field selectors, and returns the list of Harvesters that match those selectors.
func (c *FakeHarvesters) List(opts v1.ListOptions) (result *replay.HarvesterList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(harvestersResource, harvestersKind, c.ns, opts), &replay.HarvesterList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &replay.HarvesterList{}
	for _, item := range obj.(*replay.HarvesterList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested harvesters.
func (c *FakeHarvesters) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(harvestersResource, c.ns, opts))

}

// Create takes the representation of a harvester and creates it.  Returns the server's representation of the harvester, and an error, if there is any.
func (c *FakeHarvesters) Create(harvester *replay.Harvester) (result *replay.Harvester, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(harvestersResource, c.ns, harvester), &replay.Harvester{})

	if obj == nil {
		return nil, err
	}
	return obj.(*replay.Harvester), err
}

// Update takes the representation of a harvester and updates it. Returns the server's representation of the harvester, and an error, if there is any.
func (c *FakeHarvesters) Update(harvester *replay.Harvester) (result *replay.Harvester, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(harvestersResource, c.ns, harvester), &replay.Harvester{})

	if obj == nil {
		return nil, err
	}
	return obj.(*replay.Harvester), err
}

// Delete takes name of the harvester and deletes it. Returns an error if one occurs.
func (c *FakeHarvesters) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(harvestersResource, c.ns, name), &replay.Harvester{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeHarvesters) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(harvestersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &replay.HarvesterList{})
	return err
}

// Patch applies the patch and returns the patched harvester.
func (c *FakeHarvesters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *replay.Harvester, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(harvestersResource, c.ns, name, data, subresources...), &replay.Harvester{})

	if obj == nil {
		return nil, err
	}
	return obj.(*replay.Harvester), err
}
