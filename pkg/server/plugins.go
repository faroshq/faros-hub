package server

import (
	"context"
	"net/http"
	"path"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/klog/v2"

	pluginsv1alpha1 "github.com/faroshq/faros-hub/pkg/apis/plugins/v1alpha1"
)

// pluginsHandler is a http handler for plugins operations
// GET -  faros.sh/plugins - list all plugins for users
// Plugins are global, so we don't care about the user
func (s *Service) pluginsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authenticated, _, err := s.authenticator.Authenticate(r)
	if err != nil {
		klog.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !authenticated {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	// list/get
	case http.MethodGet:
		parts := strings.Split(r.URL.Path, path.Join(pathAPIVersion, pathPlugins))
		if len(parts) == 2 && parts[1] == "" { // no workspace name - list all plugins
			plugins, err := s.listPlugins(ctx)
			if err != nil {
				responsewriters.ErrorNegotiated(err, codecs, schema.GroupVersion{}, w, r)
				return
			}
			responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, pluginsv1alpha1.SchemeGroupVersion, w, r, http.StatusOK, plugins)
			return
		} else if len(parts) == 2 && parts[1] != "" { // workspace name - get workspace details
			plugin, err := s.getPlugin(ctx, strings.TrimPrefix(parts[1], "/"))
			if err != nil {
				responsewriters.ErrorNegotiated(err, codecs, schema.GroupVersion{}, w, r)
				return
			}
			responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, pluginsv1alpha1.SchemeGroupVersion, w, r, http.StatusOK, plugin)
			return
		}
	}
}

func (s *Service) listPlugins(ctx context.Context) (*pluginsv1alpha1.PluginList, error) {
	exports, err := s.kcpClient.Cluster(s.pluginsCluster).ApisV1alpha1().APIExports().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var plugins pluginsv1alpha1.PluginList
	for _, export := range exports.Items {
		// description is from apiresourceschema
		schema, err := s.kcpClient.Cluster(s.pluginsCluster).ApisV1alpha1().APIResourceSchemas().Get(ctx, export.Name, metav1.GetOptions{})
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

func (s *Service) getPlugin(ctx context.Context, name string) (*pluginsv1alpha1.Plugin, error) {
	export, err := s.kcpClient.Cluster(s.pluginsCluster).ApisV1alpha1().APIExports().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(export.Name, ".", 2)
	return &pluginsv1alpha1.Plugin{
		TypeMeta: metav1.TypeMeta{
			Kind:       pluginsv1alpha1.PluginKind,
			APIVersion: pluginsv1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: parts[1],
		},
		Spec: pluginsv1alpha1.PluginSpec{
			Version: parts[0],
		}}, nil

}
