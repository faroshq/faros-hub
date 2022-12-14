/*
Copyright 2022.

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
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMonitorings implements MonitoringInterface
type FakeMonitorings struct {
	Fake *FakePluginsV1alpha1
	ns   string
}

var monitoringsResource = schema.GroupVersionResource{Group: "plugins.faros.sh", Version: "v1alpha1", Resource: "monitorings"}

var monitoringsKind = schema.GroupVersionKind{Group: "plugins.faros.sh", Version: "v1alpha1", Kind: "Monitoring"}

// Get takes name of the monitoring, and returns the corresponding monitoring object, and an error if there is any.
func (c *FakeMonitorings) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Monitoring, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(monitoringsResource, c.ns, name), &v1alpha1.Monitoring{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Monitoring), err
}

// List takes label and field selectors, and returns the list of Monitorings that match those selectors.
func (c *FakeMonitorings) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.MonitoringList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(monitoringsResource, monitoringsKind, c.ns, opts), &v1alpha1.MonitoringList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.MonitoringList{ListMeta: obj.(*v1alpha1.MonitoringList).ListMeta}
	for _, item := range obj.(*v1alpha1.MonitoringList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested monitorings.
func (c *FakeMonitorings) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(monitoringsResource, c.ns, opts))

}

// Create takes the representation of a monitoring and creates it.  Returns the server's representation of the monitoring, and an error, if there is any.
func (c *FakeMonitorings) Create(ctx context.Context, monitoring *v1alpha1.Monitoring, opts v1.CreateOptions) (result *v1alpha1.Monitoring, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(monitoringsResource, c.ns, monitoring), &v1alpha1.Monitoring{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Monitoring), err
}

// Update takes the representation of a monitoring and updates it. Returns the server's representation of the monitoring, and an error, if there is any.
func (c *FakeMonitorings) Update(ctx context.Context, monitoring *v1alpha1.Monitoring, opts v1.UpdateOptions) (result *v1alpha1.Monitoring, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(monitoringsResource, c.ns, monitoring), &v1alpha1.Monitoring{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Monitoring), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeMonitorings) UpdateStatus(ctx context.Context, monitoring *v1alpha1.Monitoring, opts v1.UpdateOptions) (*v1alpha1.Monitoring, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(monitoringsResource, "status", c.ns, monitoring), &v1alpha1.Monitoring{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Monitoring), err
}

// Delete takes name of the monitoring and deletes it. Returns an error if one occurs.
func (c *FakeMonitorings) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(monitoringsResource, c.ns, name, opts), &v1alpha1.Monitoring{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMonitorings) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(monitoringsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.MonitoringList{})
	return err
}

// Patch applies the patch and returns the patched monitoring.
func (c *FakeMonitorings) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Monitoring, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(monitoringsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Monitoring{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Monitoring), err
}
