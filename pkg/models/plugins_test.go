package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatest(t *testing.T) {
	p := PluginsList{
		{
			Name:    "foo",
			Version: "v1",
		},
		{
			Name:    "foo",
			Version: "v2",
		},
		{
			Name:    "bar",
			Version: "v1",
		},
	}

	latest, err := p.GetLatest("foo")
	assert.NoError(t, err)
	assert.Equal(t, "v2", latest.Version)

	latest, err = p.GetLatest("bar")
	assert.NoError(t, err)
	assert.Equal(t, "v1", latest.Version)

	_, err = p.GetLatest("baz")
	assert.Error(t, err)
	assert.Equal(t, "plugin \"baz\" not found", err.Error())
}
