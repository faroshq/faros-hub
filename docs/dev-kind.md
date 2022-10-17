# Development flow with Kind

Development flow with kind is intended to test full production deployment.
It should not be used for day-to-day development. If you end-up using it
for development, something went wrong :)

## Prerequisites

- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [helm](https://helm.sh/docs/intro/install/)
- Host port 80 and 443 available
- Host resolved `local.dev.faros.sh` to `127.0.0.1`. This is done in our DNS server,
 but you might need to modify your local host configuration (like dnsMasq or `127.0.0.1 local.dev.faros.sh` in `/etc/hosts`)

1. Setup kind with `make setup-kind`
2. Deploy stack with `make deploy-kind`
