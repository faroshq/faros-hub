package models

import (
	"fmt"

	edgev1alpha1 "github.com/faroshq/faros-hub/pkg/apis/edge/v1alpha1"
)

type Plugin = edgev1alpha1.PluginSpec

// PluginsList is a list of plugins loaded via controller
type PluginsList []Plugin

// Has returns true if plugin is in the list of plugins
func (p PluginsList) Has(plugin Plugin) bool {
	for _, p := range p {
		if p.Name == plugin.Name && p.Version == plugin.Version {
			return true
		}
	}
	return false
}

// Has returns true if plugin is in the list of plugins
func (p PluginsList) GetLatest(name string) (Plugin, error) {
	latest := Plugin{
		Name:    name,
		Version: "v0",
	}
	found := false
	for _, p := range p {
		if p.Name == name && p.Version >= latest.Version {
			latest = p
			found = true
		}
	}
	if !found {
		return Plugin{}, fmt.Errorf("plugin %q not found", name)
	}
	return latest, nil
}
