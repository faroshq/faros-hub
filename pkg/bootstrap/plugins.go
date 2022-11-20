package bootstrap

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	utilplugins "github.com/faroshq/faros-hub/pkg/util/plugins"
)

func (b *bootstrap) LoadPlugins(ctx context.Context, workspace string) error {
	path := b.config.PluginsDir

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		parts := strings.Split(file.Name(), "-")
		if len(parts) != 4 {
			return fmt.Errorf("invalid plugin file name: %s", file.Name())
		}
		interfaceName := parts[1]

		p, err := utilplugins.Load(filepath.Join(path, file.Name()), interfaceName)
		if err != nil {
			return err
		}

		spew.Dump(p.Name())
		os.Exit(1)

	}

	return nil
}
