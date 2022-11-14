#!/bin/bash

hostname=$(cat ./hack/dev/kcp/values.yaml | grep externalHostname | cut -d" " -f2-)

cat << EOF > kcp.kubeconfig
apiVersion: v1
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: https://${hostname}/clusters/root
  name: kind-kcp
contexts:
- context:
    cluster: kind-kcp
    user: kind-kcp
  name: kind-kcp
current-context: kind-kcp
kind: Config
preferences: {}
users:
- name: kind-kcp
  user:
    token: admin-token
EOF


echo "Kubeconfig file created at auth-token.kubeconfig"
echo ""
echo "export KUBECONFIG=kcp.kubeconfig"
echo ""
