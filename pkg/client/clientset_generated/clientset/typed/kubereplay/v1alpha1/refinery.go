package v1alpha1

import (
	v1alpha1 "github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	scheme "github.com/lwolf/kubereplay/pkg/client/clientset_generated/clientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// RefineriesGetter has a method to return a RefineryInterface.
// A group's client should implement this interface.
type RefineriesGetter interface {
	Refineries(namespace string) RefineryInterface
}

// RefineryInterface has methods to work with Refinery resources.
type RefineryInterface interface {
	Create(*v1alpha1.Refinery) (*v1alpha1.Refinery, error)
	Update(*v1alpha1.Refinery) (*v1alpha1.Refinery, error)
	UpdateStatus(*v1alpha1.Refinery) (*v1alpha1.Refinery, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Refinery, error)
	List(opts v1.ListOptions) (*v1alpha1.RefineryList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Refinery, err error)
	RefineryExpansion
}

// refineries implements RefineryInterface
type refineries struct {
	client rest.Interface
	ns     string
}

// newRefineries returns a Refineries
func newRefineries(c *KubereplayV1alpha1Client, namespace string) *refineries {
	return &refineries{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the refinery, and returns the corresponding refinery object, and an error if there is any.
func (c *refineries) Get(name string, options v1.GetOptions) (result *v1alpha1.Refinery, err error) {
	result = &v1alpha1.Refinery{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("refineries").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Refineries that match those selectors.
func (c *refineries) List(opts v1.ListOptions) (result *v1alpha1.RefineryList, err error) {
	result = &v1alpha1.RefineryList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("refineries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested refineries.
func (c *refineries) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("refineries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a refinery and creates it.  Returns the server's representation of the refinery, and an error, if there is any.
func (c *refineries) Create(refinery *v1alpha1.Refinery) (result *v1alpha1.Refinery, err error) {
	result = &v1alpha1.Refinery{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("refineries").
		Body(refinery).
		Do().
		Into(result)
	return
}

// Update takes the representation of a refinery and updates it. Returns the server's representation of the refinery, and an error, if there is any.
func (c *refineries) Update(refinery *v1alpha1.Refinery) (result *v1alpha1.Refinery, err error) {
	result = &v1alpha1.Refinery{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("refineries").
		Name(refinery.Name).
		Body(refinery).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *refineries) UpdateStatus(refinery *v1alpha1.Refinery) (result *v1alpha1.Refinery, err error) {
	result = &v1alpha1.Refinery{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("refineries").
		Name(refinery.Name).
		SubResource("status").
		Body(refinery).
		Do().
		Into(result)
	return
}

// Delete takes name of the refinery and deletes it. Returns an error if one occurs.
func (c *refineries) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("refineries").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *refineries) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("refineries").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched refinery.
func (c *refineries) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Refinery, err error) {
	result = &v1alpha1.Refinery{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("refineries").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
