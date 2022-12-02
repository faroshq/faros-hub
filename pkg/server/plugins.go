package server

import (
	"context"
	"net/http"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/klog/v2"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	"github.com/kcp-dev/logicalcluster/v2"
)

// listPlugins lists all plugins globally
func (s *Service) listPlugins(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cluster := logicalcluster.New(s.config.ControllersPluginsWorkspace)
	client := s.kcpClient.Cluster(cluster)

	authenticated, _, err := s.authenticate(w, r)
	if err != nil || !authenticated {
		return
	}

	plugins, err := listPlugins(ctx, client)
	if err != nil {
		responsewriters.ErrorNegotiated(err, codecs, schema.GroupVersion{}, w, r)
		return
	}
	responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, pluginsv1alpha1.SchemeGroupVersion, w, r, http.StatusOK, plugins)
}

func listPlugins(ctx context.Context, client kcpclient.Interface) (*pluginsv1alpha1.PluginList, error) {
	exports, err := client.ApisV1alpha1().APIExports().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var plugins pluginsv1alpha1.PluginList
	for _, export := range exports.Items {
		// description is from apiresourceschema
		schema, err := client.ApisV1alpha1().APIResourceSchemas().Get(ctx, export.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		var description string
		for _, version := range schema.Spec.Versions {
			s, err := version.GetSchema()
			if err != nil {
				klog.Error(err)
				continue
			}
			if s.Description != "" {
				description = s.Description
				break
			}

			if description != "" {
				break
			}
		}

		parts := strings.SplitN(export.Name, ".", 2)
		plugins.Items = append(plugins.Items, pluginsv1alpha1.Plugin{
			TypeMeta: metav1.TypeMeta{
				Kind:       pluginsv1alpha1.PluginKind,
				APIVersion: pluginsv1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: parts[1],
			},
			Spec: pluginsv1alpha1.PluginSpec{
				Version:     parts[0],
				Description: description,
			},
		})
	}
	return &plugins, nil
}
