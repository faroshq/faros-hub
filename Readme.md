# KCP example repository

This repository contains a simple example of how to use the [KCP](https://github.com/kcp-dev/kcp)

## Prerequisites

Kind installed and configured.

Run `hack/dev/setup-dev.sh` to create 2 KIND clusters for services and shared compute

Start build-in KCP to bootstrap the clusters: `make run-server`

## Example internals

Once script finishes, you should have 2 KIND clusters and KCP running.
KCP should have workspace structure as bellow (`kubectl kcp workspace tree`):

```bash
.
└── root
    ├── compute
    ├── corp # - abritrary name for the organization
    │   ├── compute # - shared compute worksapce
    │   │   ├── services  # - shared compute workspace for services. Contains syncTarget to services physical cluster
    │   │   └── shared # - shared compute workspace for shared workloads. Contains syncTarget to shared physical cluster
    │   └── services
    │       └── warehouse # - services workspace for warehouse service. Runs components in services physical cluster
    └── users
        ├── user1 # - user1 workspace
        ├── user2
        └── user3
```

In addition these concepts are configured:

`Location` x 2 objects configured in `root:corp:compute:services` and `root:corp:compute:shared`
workspaces to define Location for services and shared compute clusters.

`Placement` objects configured in `root:corp:services:warehouse` workspace to define where to run
components of warehouse service (points to `root:corp:compute:services` Location)

`APIExport` objects configured in `root:corp:services:warehouse` workspace to define what API to expose
to users (`faros.sh`) API for our warehouse service.

`APIBinding` object configured for `kubernetes` api into `root:corp:services:warehouse` workspace to
be able to use kuberentes API and run workloads.

`APIBinding` object configured for `faros.sh` and `kubernetes` api into `root:users:user1` workspace to be able to use
kuberentes API and faros.sh API.

`Controller/Manager` deployed into `root:corp:services:warehouse` workspace to run components of warehouse service.
It uses 2 workspaces: `root:corp:services:warehouse` to store state in config map, and `VirtualWorkspace` for the same workspace
for `faros.sh` CRD API management.

`ServiceAccount` account in cluster `root:users:user1` to be able to use All apis.

`./dev/user1.kubeconfig` - kubeconfig file for user1 to be able to interact with all apis

## Access to the clusters

To access Root KCP cluster:
```
export KUBECONFIG=./dev/server/admin.kubeconfig`
```

To access user cluster:
```
export KUBECONFIG=./dev/user1.kubeconfig`
```

## TODO:

1. Create a RBAC example
2. Add "stop" steps before each bootstrap state for demo purposes
3. Potatoes controller does not return potatoes :/
