# Development notes

```
# setup KIND test cluster. Cluster will host external tooling for development
make setup-kind

# Run faros locally using OIDC inside the cluster. See `docs/dev-idp.md` for more details
make run-with-oidc

# login with CLI:
go run ./cmd/kubectl-faros/ login

# Once logged in you can create new workspace for the user:
go run ./cmd/kubectl-faros/ workspace create my-workspace

# If you want to observe what is happening in the cluster you can use:
# Make sure you not override current kubeconfig from command above
export KUBECONFIG=.faros/admin.kubeconfig
kubectl kcp workspace tree

# Set kubeconfig to new workspace
go run ./cmd/kubectl-faros/ workspace use my-workspace

