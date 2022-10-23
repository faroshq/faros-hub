# Controllers

Controllers in the `pkg/controllers` package are responsible for reconciling the
apis and multi-cluster objects.

There are 2 types of controllers:
- `system` - these controllers operate in single virtualWorkspace (for now) and
is responsible for tenant bootstrapping.
- `tenant` - these controllers operate in individual virtualWorkspaces (like edge.faros.sh, etc)
and is responsible for reconciling the objects in tenant workspaces.

- `system` - controllers that operate in system workspace is responsible for tenant bootstrap and system operations
- `tenants/edge` - controllers that operate at tenant level and reconciled edge (`edge.faros.sh`) APIs
- `tenants/plugins` - controllers that operate at tenant level and reconciled plugins (`plugins.faros.sh`) APIs
- `tenants/access` - controllers that operate at tenant level and reconciled access (`access.faros.sh`) APIs
