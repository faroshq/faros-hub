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
	"net/http"

	v1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	"github.com/faroshq/faros-hub/pkg/client/clientset/versioned/scheme"
	v2 "github.com/kcp-dev/logicalcluster/v2"
	rest "k8s.io/client-go/rest"
)

type PluginsV1alpha1Interface interface {
	RESTClient() rest.Interface
	ContainerRuntimesGetter
	MonitoringsGetter
	NetworksGetter
	NotificationsGetter
}

// PluginsV1alpha1Client is used to interact with features provided by the plugins.faros.sh group.
type PluginsV1alpha1Client struct {
	restClient rest.Interface
	cluster    v2.Name
}

func (c *PluginsV1alpha1Client) ContainerRuntimes(namespace string) ContainerRuntimeInterface {
	return newContainerRuntimes(c, namespace)
}

func (c *PluginsV1alpha1Client) Monitorings(namespace string) MonitoringInterface {
	return newMonitorings(c, namespace)
}

func (c *PluginsV1alpha1Client) Networks(namespace string) NetworkInterface {
	return newNetworks(c, namespace)
}

func (c *PluginsV1alpha1Client) Notifications(namespace string) NotificationInterface {
	return newNotifications(c, namespace)
}

// NewForConfig creates a new PluginsV1alpha1Client for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*PluginsV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

// NewForConfigAndClient creates a new PluginsV1alpha1Client for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*PluginsV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientForConfigAndClient(&config, h)
	if err != nil {
		return nil, err
	}
	return &PluginsV1alpha1Client{restClient: client}, nil
}

// NewForConfigOrDie creates a new PluginsV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *PluginsV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new PluginsV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *PluginsV1alpha1Client {
	return &PluginsV1alpha1Client{restClient: c}
}

// NewWithCluster creates a new PluginsV1alpha1Client for the given RESTClient and cluster.
func NewWithCluster(c rest.Interface, cluster v2.Name) *PluginsV1alpha1Client {
	return &PluginsV1alpha1Client{restClient: c, cluster: cluster}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *PluginsV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}