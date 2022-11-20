package plugins

import (
	"fmt"
	"os"
	"os/exec"

	farosplugins "github.com/faroshq/faros-hub/pkg/plugins"
	"github.com/faroshq/faros-hub/pkg/plugins/shared"
	"github.com/hashicorp/go-plugin"
)

func Load(path string) (farosplugins.Interface, error) {
	fs, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat plugin file: %s", err)
	}
	if fs.Size() == 0 {
		return nil, fmt.Errorf("%s points to an empty file", path)
	}

	// Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", path),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	//defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("plugin")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// We should have a Counter store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	p := raw.(farosplugins.Interface)

	return p, nil
}
