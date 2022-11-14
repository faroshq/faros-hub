# Development notes


Kind cluster is hosting all infrastructure required to develop locally.
It will host `kcp`, `dex`, `cert-manager`, `nginx-ingress` and `reverse-proxy` components

## Prerequisites

* `/etc/hosts` file should contain `127.0.0.1 dex.dev.faros.sh kcp.dev.faros.sh`
for local traffic routing
* ports `80` and `443` should be free for local traffic routing
* port 30443 should be free for reverse dialer so you can use remote resources
same way as it would be running inside the cluster but still be running binaries locally

## Setup

Setup KIND test cluster. It will run all required components for development and
start reverse dialer in port 30443.

```bash
make setup-kind
```

It will write hosting cluster kubeconfig:
```bash
export KUBECONFIG=.faros/admin.kubeconfig
kubectl get pods -A
```

## CLI

Login via CLI:
```bash
go run ./cmd/kubectl-faros/ login
```

Create first workspace:
```bash
go run ./cmd/kubectl-faros/ workspace create my-workspace
```

Set kubeconfig to new workspace
```bash
go run ./cmd/kubectl-faros/ workspace use my-workspace
```
