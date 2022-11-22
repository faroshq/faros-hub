package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePluginVersion(t *testing.T) {
	_, err := ParsePluginVersion("v20201231")
	assert.NoError(t, err)
	_, err = ParsePluginVersion("v2020123")
	assert.Error(t, err)
}

func TestPluginVersion_LowerThan(t *testing.T) {
	v1, _ := ParsePluginVersion("v20201231")
	v2, _ := ParsePluginVersion("v20210101")
	assert.True(t, v1.LowerThan(v2))
	assert.False(t, v2.LowerThan(v1))
}

func TestPluginVersion_HigherThan(t *testing.T) {
	v1, _ := ParsePluginVersion("v20201231")
	v2, _ := ParsePluginVersion("v20210101")
	assert.True(t, v2.HigherThan(v1))
	assert.False(t, v1.HigherThan(v2))
}
