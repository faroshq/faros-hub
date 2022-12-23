# Plugin development

Plugins code lives in separate repositories. Each repository should give a good example of how to write a plugin or how to run it locally. But when doing
e2e development, you might want to run plugin with agent code. This document will explain how to do that.

## Prerequisites

- Build plugin code inside plugin repository. This will create a plugins in `bin` directory, in example `plugins/faros-systemd-6297e18-amd64.so`.
- Set `plugins` override variable inside `agent` terminal:

```
export FAROS_PLUGINS_DIR=/path/to/agent/plugins
# example:
export FAROS_PLUGINS_DIR=/go/src/github.com/faroshq/plugin-services/plugins
```

For simplicity, you can build agent with static version of plugin. This will make
reloading easier when you are modifying plugin code. To do that, you need to set

```
# inside plugins repository
export PLUGIN_VERSION=dev
make build
```

## Agent configuration

For plugins to work, they need to be loaded into KCP backplain so APIs can
be exposed to workspaces. For this when running `controllers` or `all-in-one`
process, you need to specify


# Controller configuration & bootstrapping

Controllers will load specified plugins and expose them to workspaces.
This is done in few steps:

1. Plugins directory is being scanned for plugins. In the future this will be
replaced with plugins API and plugins will be loaded from there.

2. (TODO) Plugins are being into KCP controllers plugins workspace. And exposed from there
to workspaces. One plugin is enabled on agent spec, agent controller will expose
requested API to the workspace.

3. (TODO) Once API is available to workspace Agent controller will initiate requested
plugin instance inside workspace owned by agent. Once instance is being created
it can be used by remote agents.

4. (TODO) Server side plugin component will be initiated together with plugin loading
on the controller and will be running for multi-cluster configuration.

