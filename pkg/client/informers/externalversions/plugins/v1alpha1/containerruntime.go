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

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	versioned "github.com/faroshq/faros-hub/pkg/client/clientset/versioned"
	internalinterfaces "github.com/faroshq/faros-hub/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/faroshq/faros-hub/pkg/client/listers/plugins/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ContainerRuntimeInformer provides access to a shared informer and lister for
// ContainerRuntimes.
type ContainerRuntimeInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ContainerRuntimeLister
}

type containerRuntimeInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewContainerRuntimeInformer constructs a new informer for ContainerRuntime type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewContainerRuntimeInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredContainerRuntimeInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredContainerRuntimeInformer constructs a new informer for ContainerRuntime type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredContainerRuntimeInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return NewFilteredContainerRuntimeInformerWithOptions(client, namespace, tweakListOptions, cache.WithResyncPeriod(resyncPeriod), cache.WithIndexers(indexers))
}

func NewFilteredContainerRuntimeInformerWithOptions(client versioned.Interface, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc, opts ...cache.SharedInformerOption) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformerWithOptions(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.PluginsV1alpha1().ContainerRuntimes(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.PluginsV1alpha1().ContainerRuntimes(namespace).Watch(context.TODO(), options)
			},
		},
		&pluginsv1alpha1.ContainerRuntime{},
		opts...,
	)
}

func (f *containerRuntimeInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	indexers := cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}
	for k, v := range f.factory.ExtraNamespaceScopedIndexers() {
		indexers[k] = v
	}

	return NewFilteredContainerRuntimeInformerWithOptions(client, f.namespace,
		f.tweakListOptions,
		cache.WithResyncPeriod(resyncPeriod),
		cache.WithIndexers(indexers),
		cache.WithKeyFunction(f.factory.KeyFunction()),
	)
}

func (f *containerRuntimeInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&pluginsv1alpha1.ContainerRuntime{}, f.defaultInformer)
}

func (f *containerRuntimeInformer) Lister() v1alpha1.ContainerRuntimeLister {
	return v1alpha1.NewContainerRuntimeLister(f.Informer().GetIndexer())
}
