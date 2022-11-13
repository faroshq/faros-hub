#!/bin/bash

source .env

if [ ! -f "/usr/local/bin/kind" ]; then
 echo "Installing KIND"
 curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.14.0/kind-linux-amd64
 chmod +x ./kind
 sudo mv ./kind /usr/local/bin/kind
else
    echo "KIND already installed"
fi

CLUSTER_NAME=faros

if ! kind get clusters | grep -w -q "${CLUSTER_NAME}"; then
kind create cluster --name faros \
     --kubeconfig ./dev/faros.kubeconfig \
     --config ./hack/dev/kind/config.yaml
else
    echo "Cluster already exists"
fi

export KUBECONFIG=./dev/faros.kubeconfig

echo "Installing ingress"

#kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
# Fork of the above to add http2
kubectl apply -f https://gist.githubusercontent.com/mjudeikis/dd91434af0049378b4a24d021cceef38/raw/413600fe604bea2ccf4dcc2bd52375ebf863f35b/deploy
kubectl label nodes faros-control-plane ingress-ready="true"
kubectl label nodes faros-control-plane node-role.kubernetes.io/control-plane-

echo "Waiting for the ingress controller to become ready..."
kubectl --context "${KUBECTL_CONTEXT}" -n ingress-nginx wait --for=condition=Ready pod -l app.kubernetes.io/component=controller --timeout=5m


echo "Installing cert-manager"

helm repo add jetstack https://charts.jetstack.io
helm repo update

kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.crds.yaml
helm install \
  --wait \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.9.1


echo "Install dex"

[ ! -d "./dev/dex-chart" ] && git clone https://github.com/faroshq/dex-helm-charts -b master ./dev/dex-chart

helm upgrade -i dex ./dev/dex-chart/charts/dex \
     --values ./hack/dev/dex/values.yaml \
     --create-namespace \
     --namespace dex \
     --wait \
     --set config.connectors[0].config.clientSecret=$GITHUB_CLIENT_SECRET \
     --set config.connectors[0].config.clientID=$GITHUB_CLIENT_ID

echo "Install KCP"

# HACK to trust the dex CA
kubectl create namespace kcp
kubectl get secret dex-pki-ca -n dex -o yaml \
| sed s/"namespace: dex"/"namespace: kcp"/\
| kubectl apply -n kcp -f - | true

[ ! -d "./dev/kcp-chart" ] && git clone https://github.com/faroshq/helm-charts.git -b faros ./dev/kcp-chart

helm upgrade -i kcp ./dev/kcp-chart/charts/kcp \
     --values ./hack/dev/kcp/values.yaml \
     --namespace kcp

echo "Install Faros"

helm upgrade -i faros ./charts/faros-dev \
     --values ./hack/dev/faros/values.yaml \
     --namespace kcp

echo "Starting reverse dialer for local development"

go run ./cmd/reverse-proxy/ --clientUpstreamUrl=https://localhost:30443 client
