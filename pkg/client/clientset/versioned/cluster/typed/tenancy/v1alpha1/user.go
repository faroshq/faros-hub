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
	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v2"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"


	tenancyv1alpha1client "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/typed/tenancy/v1alpha1"
)

// UsersClusterGetter has a method to return a UserClusterInterface.
// A group's cluster client should implement this interface.
type UsersClusterGetter interface {
	Users() UserClusterInterface
}

// UserClusterInterface can operate on Users across all clusters,
// or scope down to one cluster and return a tenancyv1alpha1client.UserInterface.
type UserClusterInterface interface {
	Cluster(logicalcluster.Name) tenancyv1alpha1client.UserInterface
	List(ctx context.Context, opts metav1.ListOptions) (*tenancyv1alpha1.UserList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type usersClusterInterface struct {
	clientCache kcpclient.Cache[*tenancyv1alpha1client.TenancyV1alpha1Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *usersClusterInterface) Cluster(name logicalcluster.Name) tenancyv1alpha1client.UserInterface {
	if name == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return c.clientCache.ClusterOrDie(name).Users()
}


// List returns the entire collection of all Users across all clusters. 
func (c *usersClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*tenancyv1alpha1.UserList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).Users().List(ctx, opts)
}

// Watch begins to watch all Users across all clusters.
func (c *usersClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).Users().Watch(ctx, opts)
}