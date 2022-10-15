# Faros hub

Faros hub is Kubernetes like Edge device management control plane.
It allows to manage edge devices in the way similar to Kubernetes.

Each edge device is running reconciler which is responsible for edge device
state. Each device reads agent `spec` and updates `status` accordingly.

Project is based on [KCP](https://www.kcp.io/) allowing each tennat to be isolated
at the same time allowing centralized management.

Each tenant 'looks like' Kubernetes cluster. It has namespaces, events, rbac, etc.

Why KCP? Kubernetes is great for managing cloud resources. But it is not designed
for edge devices or multitenancy. KCP is designed for multi-tenancy and extendability.
This allows us build any platform on top of it. Use all best what kubernetes has to offer,
like client generation, rbac, consistent api, etc. And at the same time not to have
overhead of managing kubeletes, nodes, compute, etc.
It allows to write single "multi-tenant" controller and run single instance of it
managing all tenants.

![High level diagram](docs/img/hl.jph)

## Important

Project is in heavy development. It is not ready for production use.

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

- [ ] Add CLI for workspace provisioning
- [ ] Add front-proxy for JWT provider authentication
- [ ] Add image and binary builds for hub-api and edge-agent, and controllers
- [ ] Add automatic binding management for new workspaces
- [ ] Improve bootstrap (now failing on updates)
