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

package v1alpha1

import (
	"context"
	"time"

	v2 "github.com/kcp-dev/logicalcluster/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"

	v1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
	scheme "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/scheme"
)

// RegistrationsGetter has a method to return a RegistrationInterface.
// A group's client should implement this interface.
type RegistrationsGetter interface {
	Registrations(namespace string) RegistrationInterface
}

// RegistrationInterface has methods to work with Registration resources.
type RegistrationInterface interface {
	Create(ctx context.Context, registration *v1alpha1.Registration, opts v1.CreateOptions) (*v1alpha1.Registration, error)
	Update(ctx context.Context, registration *v1alpha1.Registration, opts v1.UpdateOptions) (*v1alpha1.Registration, error)
	UpdateStatus(ctx context.Context, registration *v1alpha1.Registration, opts v1.UpdateOptions) (*v1alpha1.Registration, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Registration, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.RegistrationList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Registration, err error)
	RegistrationExpansion
}

// registrations implements RegistrationInterface
type registrations struct {
	client  rest.Interface
	cluster v2.Name
	ns      string
}

// newRegistrations returns a Registrations
func newRegistrations(c *EdgeV1alpha1Client, namespace string) *registrations {
	return &registrations{
		client:  c.RESTClient(),
		cluster: c.cluster,
		ns:      namespace,
	}
}

// Get takes name of the registration, and returns the corresponding registration object, and an error if there is any.
func (c *registrations) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Registration, err error) {
	result = &v1alpha1.Registration{}
	err = c.client.Get().
		Cluster(c.cluster).
		Namespace(c.ns).
		Resource("registrations").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Registrations that match those selectors.
func (c *registrations) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.RegistrationList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.RegistrationList{}
	err = c.client.Get().
		Cluster(c.cluster).
		Namespace(c.ns).
		Resource("registrations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested registrations.
func (c *registrations) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Cluster(c.cluster).
		Namespace(c.ns).
		Resource("registrations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a registration and creates it.  Returns the server's representation of the registration, and an error, if there is any.
func (c *registrations) Create(ctx context.Context, registration *v1alpha1.Registration, opts v1.CreateOptions) (result *v1alpha1.Registration, err error) {
	result = &v1alpha1.Registration{}
	err = c.client.Post().
		Cluster(c.cluster).
		Namespace(c.ns).
		Resource("registrations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(registration).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a registration and updates it. Returns the server's representation of the registration, and an error, if there is any.
func (c *registrations) Update(ctx context.Context, registration *v1alpha1.Registration, opts v1.UpdateOptions) (result *v1alpha1.Registration, err error) {
	result = &v1alpha1.Registration{}
	err = c.client.Put().
		Cluster(c.cluster).
		Namespace(c.ns).
		Resource("registrations").
		Name(registration.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(registration).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *registrations) UpdateStatus(ctx context.Context, registration *v1alpha1.Registration, opts v1.UpdateOptions) (result *v1alpha1.Registration, err error) {
	result = &v1alpha1.Registration{}
	err = c.client.Put().
		Cluster(c.cluster).
		Namespace(c.ns).
		Resource("registrations").
		Name(registration.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(registration).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the registration and deletes it. Returns an error if one occurs.
func (c *registrations) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Cluster(c.cluster).
		Namespace(c.ns).
		Resource("registrations").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *registrations) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Cluster(c.cluster).
		Namespace(c.ns).
		Resource("registrations").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched registration.
func (c *registrations) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Registration, err error) {
	result = &v1alpha1.Registration{}
	err = c.client.Patch(pt).
		Cluster(c.cluster).
		Namespace(c.ns).
		Resource("registrations").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}