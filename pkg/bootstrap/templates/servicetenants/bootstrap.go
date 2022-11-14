package servicetenants

import (
	"context"
	"embed"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"

	utilbootstrap "github.com/faroshq/faros-hub/pkg/util/bootstrap"
)

//go:embed *.yaml
var fs embed.FS

// Bootstrap creates resources in this package by continuously retrying the list.
// This is blocking, i.e. it only returns (with error) when the context is closed or with nil when
// the bootstrapping is successfully completed.
func Bootstrap(ctx context.Context, discoveryClient discovery.DiscoveryInterface, dynamicClient dynamic.Interface, opts ...utilbootstrap.Option) error {
	return utilbootstrap.Bootstrap(ctx, discoveryClient, dynamicClient, fs, opts...)
}
