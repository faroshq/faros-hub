package plugins

import (
	"fmt"
	"os"
	goplugin "plugin"

	"github.com/faroshq/faros-hub/pkg/plugins"
)

func Load(path, name string) (plugins.Interface, error) {
	fs, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat plugin file: %s", err)
	}
	if fs.Size() == 0 {
		return nil, fmt.Errorf("%s points to an empty file", path)
	}

	plug, err := goplugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin file: %s", err)
	}

	//search for an exported Name should match the plugin name
	_p, err := plug.Lookup(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// check that loaded p is type Interface
	p, ok := _p.(plugins.Interface)
	if !ok {
		fmt.Println("The module have wrong type")
		os.Exit(-1)
	}

	return p, nil
}
