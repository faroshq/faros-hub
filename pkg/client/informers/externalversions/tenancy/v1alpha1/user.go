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
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	tenancyv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/tenancy/v1alpha1"
	versioned "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	internalinterfaces "github.com/faroshq/faros-hub/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/faroshq/faros-hub/pkg/client/listers/tenancy/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// UserInformer provides access to a shared informer and lister for
// Users.
type UserInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.UserLister
}

type userInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewUserInformer constructs a new informer for User type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewUserInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredUserInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredUserInformer constructs a new informer for User type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredUserInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return NewFilteredUserInformerWithOptions(client, tweakListOptions, cache.WithResyncPeriod(resyncPeriod), cache.WithIndexers(indexers))
}

func NewFilteredUserInformerWithOptions(client versioned.Interface, tweakListOptions internalinterfaces.TweakListOptionsFunc, opts ...cache.SharedInformerOption) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformerWithOptions(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TenancyV1alpha1().Users().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TenancyV1alpha1().Users().Watch(context.TODO(), options)
			},
		},
		&tenancyv1alpha1.User{},
		opts...,
	)
}

func (f *userInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	indexers := cache.Indexers{}
	for k, v := range f.factory.ExtraClusterScopedIndexers() {
		indexers[k] = v
	}

	return NewFilteredUserInformerWithOptions(client,
		f.tweakListOptions,
		cache.WithResyncPeriod(resyncPeriod),
		cache.WithIndexers(indexers),
		cache.WithKeyFunction(f.factory.KeyFunction()),
	)
}

func (f *userInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&tenancyv1alpha1.User{}, f.defaultInformer)
}

func (f *userInformer) Lister() v1alpha1.UserLister {
	return v1alpha1.NewUserLister(f.Informer().GetIndexer())
}
