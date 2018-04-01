package fake

import (
	v1alpha1 "github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeHarvesters implements HarvesterInterface
type FakeHarvesters struct {
	Fake *FakeKubereplayV1alpha1
	ns   string
}

var harvestersResource = schema.GroupVersionResource{Group: "kubereplay.lwolf.org", Version: "v1alpha1", Resource: "harvesters"}

var harvestersKind = schema.GroupVersionKind{Group: "kubereplay.lwolf.org", Version: "v1alpha1", Kind: "Harvester"}

// Get takes name of the harvester, and returns the corresponding harvester object, and an error if there is any.
func (c *FakeHarvesters) Get(name string, options v1.GetOptions) (result *v1alpha1.Harvester, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(harvestersResource, c.ns, name), &v1alpha1.Harvester{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Harvester), err
}

// List takes label and field selectors, and returns the list of Harvesters that match those selectors.
func (c *FakeHarvesters) List(opts v1.ListOptions) (result *v1alpha1.HarvesterList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(harvestersResource, harvestersKind, c.ns, opts), &v1alpha1.HarvesterList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.HarvesterList{}
	for _, item := range obj.(*v1alpha1.HarvesterList).Items {
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
func (c *FakeHarvesters) Create(harvester *v1alpha1.Harvester) (result *v1alpha1.Harvester, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(harvestersResource, c.ns, harvester), &v1alpha1.Harvester{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Harvester), err
}

// Update takes the representation of a harvester and updates it. Returns the server's representation of the harvester, and an error, if there is any.
func (c *FakeHarvesters) Update(harvester *v1alpha1.Harvester) (result *v1alpha1.Harvester, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(harvestersResource, c.ns, harvester), &v1alpha1.Harvester{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Harvester), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeHarvesters) UpdateStatus(harvester *v1alpha1.Harvester) (*v1alpha1.Harvester, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(harvestersResource, "status", c.ns, harvester), &v1alpha1.Harvester{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Harvester), err
}

// Delete takes name of the harvester and deletes it. Returns an error if one occurs.
func (c *FakeHarvesters) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(harvestersResource, c.ns, name), &v1alpha1.Harvester{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeHarvesters) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(harvestersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.HarvesterList{})
	return err
}

// Patch applies the patch and returns the patched harvester.
func (c *FakeHarvesters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Harvester, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(harvestersResource, c.ns, name, data, subresources...), &v1alpha1.Harvester{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Harvester), err
}
