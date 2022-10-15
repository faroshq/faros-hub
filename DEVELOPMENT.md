# Development nodes

```
# setup KIND test cluster
./hack/dev/setup-kind.sh

# create workspace and sync target for the cluster in questions

go run ./cmd/kubectl-faros/ workspace use root
go run ./cmd/kubectl-faros/ workspace create clusters --enter
go run ./cmd/kubectl-faros/ workspace create dev-env --enter
go run ./cmd/kubectl-faros/ workload sync kind --syncer-image ghcr.io/kcp-dev/kcp/syncer:main -o syncer-kind-main.yaml
```

# Deploy syncer

```
KUBECONFIG=./dev/cluster.kubeconfig kubectl apply -f "syncer-kind-main.yaml"

# Access virtual workspaces, objects

kubectl --server=https://localhost:6443/services/syncer-tunnels/clusters/root:clusters:dev-env/apis/workload.kcp.dev/v1alpha1/synctargets/kin
d/proxy get deployments -A

kubectl --server=https://localhost:6443/services/syncer-tunnels/clusters/root:clusters:dev-env/apis/workload.kcp.dev/v1alpha1/synctargets/kind/proxy get deployments -A
```
