# Controllers

Controllers in the `pkg/controllers` package are responsible for reconciling the
apis and multi-cluster objects.

- `global` - controllers that operate at global level
- `edge` - controllers that operate at tenant level and reconciled edge (`edge.faros.sh`) APIs
- `plugins` - controllers that operate at tenant level and reconciled plugins (`plugins.faros.sh`) APIs
- `access` - controllers that operate at tenant level and reconciled access (`access.faros.sh`) APIs
