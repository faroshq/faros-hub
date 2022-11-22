package version

import (
	"fmt"
	"regexp"
)

var PluginVersionRegex = regexp.MustCompile("v20[0-9]{6}")

type PluginVersion string

func ParsePluginVersion(version string) (*PluginVersion, error) {
	if !PluginVersionRegex.MatchString(version) {
		return nil, fmt.Errorf("invalid plugin version %q", version)
	}
	v := PluginVersion(version)
	return &v, nil
}

func (v *PluginVersion) LowerThan(other *PluginVersion) bool {
	return v.String() < other.String()
}

func (v *PluginVersion) HigherThan(other *PluginVersion) bool {
	return v.String() > other.String()
}

func (v *PluginVersion) String() string {
	return string(*v)
}
