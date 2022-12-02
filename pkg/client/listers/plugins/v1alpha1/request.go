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
	kcpcache "github.com/kcp-dev/apimachinery/pkg/cache"	
	"github.com/kcp-dev/logicalcluster/v2"
	
	"k8s.io/client-go/tools/cache"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/api/errors"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	)

// RequestClusterLister can list Requests across all workspaces, or scope down to a RequestLister for one workspace.
// All objects returned here must be treated as read-only.
type RequestClusterLister interface {
	// List lists all Requests in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*pluginsv1alpha1.Request, err error)
	// Cluster returns a lister that can list and get Requests in one workspace.
Cluster(cluster logicalcluster.Name) RequestLister
RequestClusterListerExpansion
}

type requestClusterLister struct {
	indexer cache.Indexer
}

// NewRequestClusterLister returns a new RequestClusterLister.
// We assume that the indexer:
// - is fed by a cross-workspace LIST+WATCH
// - uses kcpcache.MetaClusterNamespaceKeyFunc as the key function
// - has the kcpcache.ClusterIndex as an index
func NewRequestClusterLister(indexer cache.Indexer) *requestClusterLister {
	return &requestClusterLister{indexer: indexer}
}

// List lists all Requests in the indexer across all workspaces.
func (s *requestClusterLister) List(selector labels.Selector) (ret []*pluginsv1alpha1.Request, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*pluginsv1alpha1.Request))
	})
	return ret, err
}

// Cluster scopes the lister to one workspace, allowing users to list and get Requests.
func (s *requestClusterLister) Cluster(cluster logicalcluster.Name) RequestLister {
return &requestLister{indexer: s.indexer, cluster: cluster}
}

// RequestLister can list all Requests, or get one in particular.
// All objects returned here must be treated as read-only.
type RequestLister interface {
	// List lists all Requests in the workspace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*pluginsv1alpha1.Request, err error)
// Get retrieves the Request from the indexer for a given workspace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*pluginsv1alpha1.Request, error)
RequestListerExpansion
}
// requestLister can list all Requests inside a workspace.
type requestLister struct {
	indexer cache.Indexer
	cluster logicalcluster.Name
}

// List lists all Requests in the indexer for a workspace.
func (s *requestLister) List(selector labels.Selector) (ret []*pluginsv1alpha1.Request, err error) {
	err = kcpcache.ListAllByCluster(s.indexer, s.cluster, selector, func(i interface{}) {
		ret = append(ret, i.(*pluginsv1alpha1.Request))
	})
	return ret, err
}

// Get retrieves the Request from the indexer for a given workspace and name.
func (s *requestLister) Get(name string) (*pluginsv1alpha1.Request, error) {
	key := kcpcache.ToClusterAwareKey(s.cluster.String(), "", name)
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(pluginsv1alpha1.Resource("Request"), name)
	}
	return obj.(*pluginsv1alpha1.Request), nil
}
// NewRequestLister returns a new RequestLister.
// We assume that the indexer:
// - is fed by a workspace-scoped LIST+WATCH
// - uses cache.MetaNamespaceKeyFunc as the key function
func NewRequestLister(indexer cache.Indexer) *requestScopedLister {
	return &requestScopedLister{indexer: indexer}
}

// requestScopedLister can list all Requests inside a workspace.
type requestScopedLister struct {
	indexer cache.Indexer
}

// List lists all Requests in the indexer for a workspace.
func (s *requestScopedLister) List(selector labels.Selector) (ret []*pluginsv1alpha1.Request, err error) {
	err = cache.ListAll(s.indexer, selector, func(i interface{}) {
		ret = append(ret, i.(*pluginsv1alpha1.Request))
	})
	return ret, err
}

// Get retrieves the Request from the indexer for a given workspace and name.
func (s *requestScopedLister) Get(name string) (*pluginsv1alpha1.Request, error) {
	key := name
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(pluginsv1alpha1.Resource("Request"), name)
	}
	return obj.(*pluginsv1alpha1.Request), nil
}
