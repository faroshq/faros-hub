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

package versioned

import (
	"fmt"
	"net/http"

	accessv1alpha1 "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/typed/access/v1alpha1"
	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/typed/edge/v1alpha1"
	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/client/clientset/versioned/typed/plugins/v1alpha1"
	v2 "github.com/kcp-dev/logicalcluster/v2"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type ClusterInterface interface {
	Cluster(name v2.Name) Interface
}

type Cluster struct {
	*scopedClientset
}

// Cluster sets the cluster for a Clientset.
func (c *Cluster) Cluster(name v2.Name) Interface {
	return &Clientset{
		scopedClientset: c.scopedClientset,
		cluster:         name,
	}
}

// NewClusterForConfig creates a new Cluster for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewClusterForConfig will generate a rate-limiter in configShallowCopy.
func NewClusterForConfig(c *rest.Config) (*Cluster, error) {
	cs, err := NewForConfig(c)
	if err != nil {
		return nil, err
	}
	return &Cluster{scopedClientset: cs.scopedClientset}, nil
}

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	AccessV1alpha1() accessv1alpha1.AccessV1alpha1Interface
	EdgeV1alpha1() edgev1alpha1.EdgeV1alpha1Interface
	PluginsV1alpha1() pluginsv1alpha1.PluginsV1alpha1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*scopedClientset
	cluster v2.Name
}

// scopedClientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type scopedClientset struct {
	*discovery.DiscoveryClient
	accessV1alpha1  *accessv1alpha1.AccessV1alpha1Client
	edgeV1alpha1    *edgev1alpha1.EdgeV1alpha1Client
	pluginsV1alpha1 *pluginsv1alpha1.PluginsV1alpha1Client
}

// AccessV1alpha1 retrieves the AccessV1alpha1Client
func (c *Clientset) AccessV1alpha1() accessv1alpha1.AccessV1alpha1Interface {
	return accessv1alpha1.NewWithCluster(c.accessV1alpha1.RESTClient(), c.cluster)
}

// EdgeV1alpha1 retrieves the EdgeV1alpha1Client
func (c *Clientset) EdgeV1alpha1() edgev1alpha1.EdgeV1alpha1Interface {
	return edgev1alpha1.NewWithCluster(c.edgeV1alpha1.RESTClient(), c.cluster)
}

// PluginsV1alpha1 retrieves the PluginsV1alpha1Client
func (c *Clientset) PluginsV1alpha1() pluginsv1alpha1.PluginsV1alpha1Interface {
	return pluginsv1alpha1.NewWithCluster(c.pluginsV1alpha1.RESTClient(), c.cluster)
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient.WithCluster(c.cluster)
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c

	if configShallowCopy.UserAgent == "" {
		configShallowCopy.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	// share the transport between all clients
	httpClient, err := rest.HTTPClientFor(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	return NewForConfigAndClient(&configShallowCopy, httpClient)
}

// NewForConfigAndClient creates a new Clientset for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfigAndClient will generate a rate-limiter in configShallowCopy.
func NewForConfigAndClient(c *rest.Config, httpClient *http.Client) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}

	var cs scopedClientset
	var err error
	cs.accessV1alpha1, err = accessv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.edgeV1alpha1, err = edgev1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	cs.pluginsV1alpha1, err = pluginsv1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	return &Clientset{scopedClientset: &cs}, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	cs, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs scopedClientset
	cs.accessV1alpha1 = accessv1alpha1.New(c)
	cs.edgeV1alpha1 = edgev1alpha1.New(c)
	cs.pluginsV1alpha1 = pluginsv1alpha1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &Clientset{scopedClientset: &cs}
}
