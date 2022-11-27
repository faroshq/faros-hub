# Controllers

Controllers in the `pkg/controllers` package are responsible for reconciling the
apis and multi-cluster objects.

There are 2 types of controllers:
- `service` - these controllers operate in single virtualWorkspace (for now) and
is responsible for tenant bootstrap. Like workspaces, users, roles for workspaces.
- `tenants` - these controllers operate in individual virtualWorkspaces (like edge.faros.sh, etc)
and is responsible for reconciling the objects in tenant workspaces.
