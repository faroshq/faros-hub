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
	kcptesting "github.com/kcp-dev/client-go/third_party/k8s.io/client-go/testing"
	"github.com/kcp-dev/logicalcluster/v2"

	"k8s.io/client-go/rest"
	kcptenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/cluster/typed/tenancy/v1alpha1"
	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/typed/tenancy/v1alpha1"
)

var _ kcptenancyv1alpha1.TenancyV1alpha1ClusterInterface = (*TenancyV1alpha1ClusterClient)(nil)

type TenancyV1alpha1ClusterClient struct {
	*kcptesting.Fake 
}

func (c *TenancyV1alpha1ClusterClient) Cluster(cluster logicalcluster.Name) tenancyv1alpha1.TenancyV1alpha1Interface {
	if cluster == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}
	return &TenancyV1alpha1Client{Fake: c.Fake, Cluster: cluster}
}


func (c *TenancyV1alpha1ClusterClient) Users() kcptenancyv1alpha1.UserClusterInterface {
	return &usersClusterClient{Fake: c.Fake}
}

func (c *TenancyV1alpha1ClusterClient) Workspaces() kcptenancyv1alpha1.WorkspaceClusterInterface {
	return &workspacesClusterClient{Fake: c.Fake}
}
var _ tenancyv1alpha1.TenancyV1alpha1Interface = (*TenancyV1alpha1Client)(nil)

type TenancyV1alpha1Client struct {
	*kcptesting.Fake
	Cluster logicalcluster.Name
}

func (c *TenancyV1alpha1Client) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}


func (c *TenancyV1alpha1Client) Users() tenancyv1alpha1.UserInterface {
	return &usersClient{Fake: c.Fake, Cluster: c.Cluster}
}

func (c *TenancyV1alpha1Client) Workspaces() tenancyv1alpha1.WorkspaceInterface {
	return &workspacesClient{Fake: c.Fake, Cluster: c.Cluster}
}
