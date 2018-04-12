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
*/package v1alpha1

import (
	v1alpha1 "github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	scheme "github.com/lwolf/kubereplay/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// HarvestersGetter has a method to return a HarvesterInterface.
// A group's client should implement this interface.
type HarvestersGetter interface {
	Harvesters(namespace string) HarvesterInterface
}

// HarvesterInterface has methods to work with Harvester resources.
type HarvesterInterface interface {
	Create(*v1alpha1.Harvester) (*v1alpha1.Harvester, error)
	Update(*v1alpha1.Harvester) (*v1alpha1.Harvester, error)
	UpdateStatus(*v1alpha1.Harvester) (*v1alpha1.Harvester, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Harvester, error)
	List(opts v1.ListOptions) (*v1alpha1.HarvesterList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Harvester, err error)
	HarvesterExpansion
}

// harvesters implements HarvesterInterface
type harvesters struct {
	client rest.Interface
	ns     string
}

// newHarvesters returns a Harvesters
func newHarvesters(c *KubereplayV1alpha1Client, namespace string) *harvesters {
	return &harvesters{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the harvester, and returns the corresponding harvester object, and an error if there is any.
func (c *harvesters) Get(name string, options v1.GetOptions) (result *v1alpha1.Harvester, err error) {
	result = &v1alpha1.Harvester{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("harvesters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Harvesters that match those selectors.
func (c *harvesters) List(opts v1.ListOptions) (result *v1alpha1.HarvesterList, err error) {
	result = &v1alpha1.HarvesterList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("harvesters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested harvesters.
func (c *harvesters) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("harvesters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a harvester and creates it.  Returns the server's representation of the harvester, and an error, if there is any.
func (c *harvesters) Create(harvester *v1alpha1.Harvester) (result *v1alpha1.Harvester, err error) {
	result = &v1alpha1.Harvester{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("harvesters").
		Body(harvester).
		Do().
		Into(result)
	return
}

// Update takes the representation of a harvester and updates it. Returns the server's representation of the harvester, and an error, if there is any.
func (c *harvesters) Update(harvester *v1alpha1.Harvester) (result *v1alpha1.Harvester, err error) {
	result = &v1alpha1.Harvester{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("harvesters").
		Name(harvester.Name).
		Body(harvester).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *harvesters) UpdateStatus(harvester *v1alpha1.Harvester) (result *v1alpha1.Harvester, err error) {
	result = &v1alpha1.Harvester{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("harvesters").
		Name(harvester.Name).
		SubResource("status").
		Body(harvester).
		Do().
		Into(result)
	return
}

// Delete takes name of the harvester and deletes it. Returns an error if one occurs.
func (c *harvesters) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("harvesters").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *harvesters) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("harvesters").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched harvester.
func (c *harvesters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Harvester, err error) {
	result = &v1alpha1.Harvester{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("harvesters").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
