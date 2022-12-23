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
//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by kcp code-generator. DO NOT EDIT.

package v1alpha1

import (
	"github.com/kcp-dev/logicalcluster/v3"
	kcptesting "github.com/kcp-dev/client-go/third_party/k8s.io/client-go/testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"

	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/testing"

	"k8s.io/apimachinery/pkg/types"

	kcppluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/cluster/typed/plugins/v1alpha1"

	pluginsv1alpha1client "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/typed/plugins/v1alpha1"
)

var bindingsResource = schema.GroupVersionResource{Group: "plugins.faros.sh", Version: "v1alpha1", Resource: "bindings"}
var bindingsKind = schema.GroupVersionKind{Group: "plugins.faros.sh", Version: "v1alpha1", Kind: "Binding"}

type bindingsClusterClient struct {
	*kcptesting.Fake
}

// Cluster scopes the client down to a particular cluster.
func (c *bindingsClusterClient) Cluster(clusterPath logicalcluster.Path) kcppluginsv1alpha1.BindingsNamespacer {
	if clusterPath == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return &bindingsNamespacer{Fake: c.Fake, ClusterPath: clusterPath}
}


// List takes label and field selectors, and returns the list of Bindings that match those selectors across all clusters.
func (c *bindingsClusterClient) List(ctx context.Context, opts metav1.ListOptions) (*pluginsv1alpha1.BindingList, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewListAction(bindingsResource, bindingsKind, logicalcluster.Wildcard, metav1.NamespaceAll, opts), &pluginsv1alpha1.BindingList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &pluginsv1alpha1.BindingList{ListMeta: obj.(*pluginsv1alpha1.BindingList).ListMeta}
	for _, item := range obj.(*pluginsv1alpha1.BindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested Bindings across all clusters.
func (c *bindingsClusterClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(kcptesting.NewWatchAction(bindingsResource, logicalcluster.Wildcard, metav1.NamespaceAll, opts))
}
type bindingsNamespacer struct {
	*kcptesting.Fake
	ClusterPath logicalcluster.Path
}

func (n *bindingsNamespacer) Namespace(namespace string) pluginsv1alpha1client.BindingInterface {
	return &bindingsClient{Fake: n.Fake, ClusterPath: n.ClusterPath, Namespace: namespace}
}
type bindingsClient struct {
	*kcptesting.Fake
	ClusterPath logicalcluster.Path
	Namespace string
}


func (c *bindingsClient) Create(ctx context.Context, binding *pluginsv1alpha1.Binding, opts metav1.CreateOptions) (*pluginsv1alpha1.Binding, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewCreateAction(bindingsResource, c.ClusterPath, c.Namespace, binding), &pluginsv1alpha1.Binding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*pluginsv1alpha1.Binding), err
}

func (c *bindingsClient) Update(ctx context.Context, binding *pluginsv1alpha1.Binding, opts metav1.UpdateOptions) (*pluginsv1alpha1.Binding, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewUpdateAction(bindingsResource, c.ClusterPath, c.Namespace, binding), &pluginsv1alpha1.Binding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*pluginsv1alpha1.Binding), err
}

func (c *bindingsClient) UpdateStatus(ctx context.Context, binding *pluginsv1alpha1.Binding, opts metav1.UpdateOptions) (*pluginsv1alpha1.Binding, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewUpdateSubresourceAction(bindingsResource, c.ClusterPath, "status", c.Namespace, binding), &pluginsv1alpha1.Binding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*pluginsv1alpha1.Binding), err
}

func (c *bindingsClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.Invokes(kcptesting.NewDeleteActionWithOptions(bindingsResource, c.ClusterPath, c.Namespace, name, opts), &pluginsv1alpha1.Binding{})
	return err
}

func (c *bindingsClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := kcptesting.NewDeleteCollectionAction(bindingsResource, c.ClusterPath, c.Namespace, listOpts)

	_, err := c.Fake.Invokes(action, &pluginsv1alpha1.BindingList{})
	return err
}

func (c *bindingsClient) Get(ctx context.Context, name string, options metav1.GetOptions) (*pluginsv1alpha1.Binding, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewGetAction(bindingsResource, c.ClusterPath, c.Namespace, name), &pluginsv1alpha1.Binding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*pluginsv1alpha1.Binding), err
}

// List takes label and field selectors, and returns the list of Bindings that match those selectors.
func (c *bindingsClient) List(ctx context.Context, opts metav1.ListOptions) (*pluginsv1alpha1.BindingList, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewListAction(bindingsResource, bindingsKind, c.ClusterPath, c.Namespace, opts), &pluginsv1alpha1.BindingList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &pluginsv1alpha1.BindingList{ListMeta: obj.(*pluginsv1alpha1.BindingList).ListMeta}
	for _, item := range obj.(*pluginsv1alpha1.BindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

func (c *bindingsClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(kcptesting.NewWatchAction(bindingsResource, c.ClusterPath, c.Namespace, opts))
}

func (c *bindingsClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*pluginsv1alpha1.Binding, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewPatchSubresourceAction(bindingsResource, c.ClusterPath, c.Namespace, name, pt, data, subresources...), &pluginsv1alpha1.Binding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*pluginsv1alpha1.Binding), err
}