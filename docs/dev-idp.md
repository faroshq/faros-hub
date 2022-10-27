# Dev IDP

For dev we use Dex as identity provider. More of it: https://github.com/dexidp/dex.git

If you setup kind cluster with `make setup-kind` it will install Dex and configure it to use static users.

You need to make sure `dex.dev.faros.sh` resolved to 127.0.0.1 in your machine.
Via `/etc/hosts` or other means.

Make sure your .env file has development application id and secret. See env.example for more details.

1. Start faros local process with IDP provider:
`hack/run-with-oidc.sh`

2. Start login app in separate terminal:
```bash
cd hack/dev/dex/app
go install .
app --issuer https://dex.dev.faros.sh --issuer-root-ca ../ssl/ca.pem --client-id faros
```

Open prompted page and login without entering any details.
Once successfully logged in you get Token and Refresh Token values.

To test these in Faros:

k kcp workspace use root:faros-system:tenants

```bash
cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: workspace-admin
  clusterName: root:faros
rules:
- apiGroups:
  - tenancy.kcp.dev
  resources:
  - workspaces/content
  resourceNames:
  - controllers
  verbs:
  - admin
EOF
```

```bash
cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cluster-admin
  clusterName: root:faros
subjects:
- kind: User
  name: faros-sso:mangirdas@judeikis.lt
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: faros-system-workspace-admin
EOF

```
export TOKEN=<ID token from login app>
curl -H "Authorization: Bearer $TOKEN" -k https://192.168.1.138:6443/clusters/root:faros-system:controllers
