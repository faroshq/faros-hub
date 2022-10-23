# Faros

Faros hub is Kubernetes native Edge device management control plane.
It allows to manage edge devices in the way similar to Kubernetes.

Faros hub will enable you to connect remote device to central hub and deploy plugins to it.

## Getting started

Start by installing Faros hub on your Kubernetes cluster.

```bash
TBC
```

Create first workspace and register your first device.

```bash
TBC
```

Once device registered you can deploy plugins to it.

## Plugins

Plugins are go binaries containing Kubernetes reconciler. Faros will distribute
plugins to edge devices and run them on the device with right context. At the same time
exposing Plugins API into workspaces, where plugins are enabled.

Plugins will be able to communicate with Faros hub and other plugins by interacting
with Kubernetes API.

Example for network plugin in your workspace:
```bash
apiVersion: plugins.faros.sh/v1alpha1
kind: Network
metadata:
  name: wireguard
spec:
  type: wireguard
  config:
    privateKey: "TBC"
    listenPort: 51820
    peers:
      - publicKey: "TBC"
        allowedIPs: "
```

Example for docker plugin in your workspace:
```bash
apiVersion: plugins.faros.sh/v1alpha1
kind: ContainerRuntime
metadata:
  name: docker-runtime
spec:
  type: docker
  config: TBC
```

Configure your device to use wireguard network and docker runtime plugins.

```bash
apiVersion: edge.faros.sh/v1alpha1
kind: Agent
metadata:
  name: foo1
spec:
  plugins:
   - name: wireguard
     config: TBC
   - name: docker-runtime
     config: TBC
```

This will deploy wireguard network plugin to your device and configure it to use it.
At the same time it will deploy docker runtime plugin and configure it to use it.

You can write server and client side plugins. In example you can write plugin that
reconciles and run in your edge device, and server side plugin that will be able to
manage data coming from your device.

Each device reads agent `spec` and updates `status` accordingly.

Anybody can write plugins and publish to faros marketplace.

![High level diagram](docs/img/hl.jpg)

## Important

Project is in heavy development. It is not ready for production use.
Evolving architecture overview can be found here:
https://miro.com/app/board/o9J_lob-CMw=/?share_link_id=60410961009

## Getting started

Currently project still needs [`kubectl-kcp`](https://github.com/kcp-dev/kcp) to be installed. It will be replaced
by `kubectl-faros` in the future with opinionated configuration.

Run Faros hub "all-in-one" configuration:

```bash
go run ./cmd/hub-api start --all-in-one
```

This will start hub-api, reconciler.

Create first workpace/virtual cluster:

```bash
export KUBECONFIG=.faros/admin.kubeconfig
kubectl kcp workspace use root
kubectl kcp workspace create tenant1 --enter
```

Controller will run in tenant `root:compute:controllers` so for now this name
is reserved. It will be changed in the future.

Create APIBinding for tenant1. This will be automated in the future.
This allows faros controller to expose Faros API into virtual cluster without
need to run it there.

```bash
kubectl create -f config/samples/binding.yaml
```

Create first agent:

```bash
go run ./cmd/kubectl-faros agent generate agent1 -o agent1.kubeconfig
```

This will create `Registration` object in `root:tenant1` namespace.
`Registration` object is used to register agent with hub. It is backed by `serviceAccount`,
`Role`, `RoleBinding` and `Secret` objects.

Same registration object can be used to register multiple agents.

Open new terminal and run agent with generated kubeconfig:

export FAROS_AGENT_NAME=agent1
export FAROS_AGENT_NAMESPACE=default # default kubernetes namespace
export KUBECONFIG=agent1.kubeconfig

go run ./cmd/edge-agent
```

In the first terminal you should see Agent reporting to hub:

```bash
[mjudeikis@unknown faros-hub]$ kubectl get agent agent -o yaml
apiVersion: edge.faros.sh/v1alpha1
kind: Agent
metadata:
  annotations:
    kcp.dev/cluster: root:compute:clusters
  creationTimestamp: "2022-10-15T11:17:11Z"
  generation: 1
  name: agent
  namespace: default
  resourceVersion: "15018"
  uid: 18c9b323-dc48-4e67-81bd-aade4da83efa
spec: {}
status:
  conditions:
  - lastTransitionTime: "2022-10-15T12:06:39Z"
    status: "True"
    type: Ready
```


# Roadmap

- [x] Add plugins API
- [ ] Add plugins runner (go plugins) to edge-agent
- [ ] Add CLI for workspace provisioning
- [ ] Add OIDC provider example
- [ ] Add user management pattern (system workspace) with binding provisioning
- [ ] Add image and binary builds for hub-api and edge-agent, and controllers
- [x] Add automatic binding management for new workspaces
- [ ] Improve bootstrap (now failing on updates)
- [ ] Add HA deployment pattern (helm-chart)
