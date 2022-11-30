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
	"github.com/kcp-dev/logicalcluster/v2"
	kcptesting "github.com/kcp-dev/client-go/third_party/k8s.io/client-go/testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"

	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/testing"

	"k8s.io/apimachinery/pkg/types"

	tenancyv1alpha1client "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/typed/tenancy/v1alpha1"
)

var workspacesResource = schema.GroupVersionResource{Group: "tenancy.faros.sh", Version: "v1alpha1", Resource: "workspaces"}
var workspacesKind = schema.GroupVersionKind{Group: "tenancy.faros.sh", Version: "v1alpha1", Kind: "Workspace"}

type workspacesClusterClient struct {
	*kcptesting.Fake
}

// Cluster scopes the client down to a particular cluster.
func (c *workspacesClusterClient) Cluster(cluster logicalcluster.Name) tenancyv1alpha1client.WorkspaceInterface {
	if cluster == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return &workspacesClient{Fake: c.Fake, Cluster: cluster}
}


// List takes label and field selectors, and returns the list of Workspaces that match those selectors across all clusters.
func (c *workspacesClusterClient) List(ctx context.Context, opts metav1.ListOptions) (*tenancyv1alpha1.WorkspaceList, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootListAction(workspacesResource, workspacesKind, logicalcluster.Wildcard, opts), &tenancyv1alpha1.WorkspaceList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &tenancyv1alpha1.WorkspaceList{ListMeta: obj.(*tenancyv1alpha1.WorkspaceList).ListMeta}
	for _, item := range obj.(*tenancyv1alpha1.WorkspaceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested Workspaces across all clusters.
func (c *workspacesClusterClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(kcptesting.NewRootWatchAction(workspacesResource, logicalcluster.Wildcard, opts))
}
type workspacesClient struct {
	*kcptesting.Fake
	Cluster logicalcluster.Name
	
}


func (c *workspacesClient) Create(ctx context.Context, workspace *tenancyv1alpha1.Workspace, opts metav1.CreateOptions) (*tenancyv1alpha1.Workspace, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootCreateAction(workspacesResource, c.Cluster, workspace), &tenancyv1alpha1.Workspace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*tenancyv1alpha1.Workspace), err
}

func (c *workspacesClient) Update(ctx context.Context, workspace *tenancyv1alpha1.Workspace, opts metav1.UpdateOptions) (*tenancyv1alpha1.Workspace, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootUpdateAction(workspacesResource, c.Cluster, workspace), &tenancyv1alpha1.Workspace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*tenancyv1alpha1.Workspace), err
}

func (c *workspacesClient) UpdateStatus(ctx context.Context, workspace *tenancyv1alpha1.Workspace, opts metav1.UpdateOptions) (*tenancyv1alpha1.Workspace, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootUpdateSubresourceAction(workspacesResource, c.Cluster, "status", workspace), &tenancyv1alpha1.Workspace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*tenancyv1alpha1.Workspace), err
}

func (c *workspacesClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.Invokes(kcptesting.NewRootDeleteActionWithOptions(workspacesResource, c.Cluster, name, opts), &tenancyv1alpha1.Workspace{})
	return err
}

func (c *workspacesClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := kcptesting.NewRootDeleteCollectionAction(workspacesResource, c.Cluster, listOpts)

	_, err := c.Fake.Invokes(action, &tenancyv1alpha1.WorkspaceList{})
	return err
}

func (c *workspacesClient) Get(ctx context.Context, name string, options metav1.GetOptions) (*tenancyv1alpha1.Workspace, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootGetAction(workspacesResource, c.Cluster, name), &tenancyv1alpha1.Workspace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*tenancyv1alpha1.Workspace), err
}

// List takes label and field selectors, and returns the list of Workspaces that match those selectors.
func (c *workspacesClient) List(ctx context.Context, opts metav1.ListOptions) (*tenancyv1alpha1.WorkspaceList, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootListAction(workspacesResource, workspacesKind, c.Cluster, opts), &tenancyv1alpha1.WorkspaceList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &tenancyv1alpha1.WorkspaceList{ListMeta: obj.(*tenancyv1alpha1.WorkspaceList).ListMeta}
	for _, item := range obj.(*tenancyv1alpha1.WorkspaceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

func (c *workspacesClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(kcptesting.NewRootWatchAction(workspacesResource, c.Cluster, opts))
}

func (c *workspacesClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*tenancyv1alpha1.Workspace, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewRootPatchSubresourceAction(workspacesResource, c.Cluster, name, pt, data, subresources...), &tenancyv1alpha1.Workspace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*tenancyv1alpha1.Workspace), err
}
