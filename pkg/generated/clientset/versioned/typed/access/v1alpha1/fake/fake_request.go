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

	v1alpha1 "github.com/faroshq/faros-hub/pkg/apis/access/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRequests implements RequestInterface
type FakeRequests struct {
	Fake *FakeAccessV1alpha1
	ns   string
}

var requestsResource = schema.GroupVersionResource{Group: "access", Version: "v1alpha1", Resource: "requests"}

var requestsKind = schema.GroupVersionKind{Group: "access", Version: "v1alpha1", Kind: "Request"}

// Get takes name of the request, and returns the corresponding request object, and an error if there is any.
func (c *FakeRequests) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Request, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(requestsResource, c.ns, name), &v1alpha1.Request{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Request), err
}

// List takes label and field selectors, and returns the list of Requests that match those selectors.
func (c *FakeRequests) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.RequestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(requestsResource, requestsKind, c.ns, opts), &v1alpha1.RequestList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.RequestList{ListMeta: obj.(*v1alpha1.RequestList).ListMeta}
	for _, item := range obj.(*v1alpha1.RequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested requests.
func (c *FakeRequests) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(requestsResource, c.ns, opts))

}

// Create takes the representation of a request and creates it.  Returns the server's representation of the request, and an error, if there is any.
func (c *FakeRequests) Create(ctx context.Context, request *v1alpha1.Request, opts v1.CreateOptions) (result *v1alpha1.Request, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(requestsResource, c.ns, request), &v1alpha1.Request{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Request), err
}

// Update takes the representation of a request and updates it. Returns the server's representation of the request, and an error, if there is any.
func (c *FakeRequests) Update(ctx context.Context, request *v1alpha1.Request, opts v1.UpdateOptions) (result *v1alpha1.Request, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(requestsResource, c.ns, request), &v1alpha1.Request{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Request), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRequests) UpdateStatus(ctx context.Context, request *v1alpha1.Request, opts v1.UpdateOptions) (*v1alpha1.Request, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(requestsResource, "status", c.ns, request), &v1alpha1.Request{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Request), err
}

// Delete takes name of the request and deletes it. Returns an error if one occurs.
func (c *FakeRequests) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(requestsResource, c.ns, name, opts), &v1alpha1.Request{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRequests) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(requestsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.RequestList{})
	return err
}

// Patch applies the patch and returns the patched request.
func (c *FakeRequests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Request, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(requestsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Request{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Request), err
}
