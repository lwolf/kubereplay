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
package v1alpha1

import (
	v1alpha1 "github.com/lwolf/kube-replay/pkg/apis/replay/v1alpha1"
	scheme "github.com/lwolf/kube-replay/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SilosGetter has a method to return a SiloInterface.
// A group's client should implement this interface.
type SilosGetter interface {
	Silos(namespace string) SiloInterface
}

// SiloInterface has methods to work with Silo resources.
type SiloInterface interface {
	Create(*v1alpha1.Silo) (*v1alpha1.Silo, error)
	Update(*v1alpha1.Silo) (*v1alpha1.Silo, error)
	UpdateStatus(*v1alpha1.Silo) (*v1alpha1.Silo, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Silo, error)
	List(opts v1.ListOptions) (*v1alpha1.SiloList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Silo, err error)
	SiloExpansion
}

// silos implements SiloInterface
type silos struct {
	client rest.Interface
	ns     string
}

// newSilos returns a Silos
func newSilos(c *KubereplayV1alpha1Client, namespace string) *silos {
	return &silos{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the silo, and returns the corresponding silo object, and an error if there is any.
func (c *silos) Get(name string, options v1.GetOptions) (result *v1alpha1.Silo, err error) {
	result = &v1alpha1.Silo{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("silos").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Silos that match those selectors.
func (c *silos) List(opts v1.ListOptions) (result *v1alpha1.SiloList, err error) {
	result = &v1alpha1.SiloList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("silos").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested silos.
func (c *silos) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("silos").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a silo and creates it.  Returns the server's representation of the silo, and an error, if there is any.
func (c *silos) Create(silo *v1alpha1.Silo) (result *v1alpha1.Silo, err error) {
	result = &v1alpha1.Silo{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("silos").
		Body(silo).
		Do().
		Into(result)
	return
}

// Update takes the representation of a silo and updates it. Returns the server's representation of the silo, and an error, if there is any.
func (c *silos) Update(silo *v1alpha1.Silo) (result *v1alpha1.Silo, err error) {
	result = &v1alpha1.Silo{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("silos").
		Name(silo.Name).
		Body(silo).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *silos) UpdateStatus(silo *v1alpha1.Silo) (result *v1alpha1.Silo, err error) {
	result = &v1alpha1.Silo{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("silos").
		Name(silo.Name).
		SubResource("status").
		Body(silo).
		Do().
		Into(result)
	return
}

// Delete takes name of the silo and deletes it. Returns an error if one occurs.
func (c *silos) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("silos").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *silos) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("silos").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched silo.
func (c *silos) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Silo, err error) {
	result = &v1alpha1.Silo{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("silos").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
