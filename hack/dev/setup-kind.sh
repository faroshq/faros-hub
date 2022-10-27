#!/bin/bash


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
     --config ./hack/dev/config.yaml
else
    echo "Cluster already exists"
fi

export KUBECONFIG=./dev/faros.kubeconfig

echo "Installing ingress"

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
kubectl label nodes faros-control-plane ingress-ready="true"
kubectl label nodes faros-control-plane node-role.kubernetes.io/control-plane-

echo "Installing cert-manager"

helm repo add jetstack https://charts.jetstack.io
helm repo update

kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.crds.yaml
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.9.1


echo "Install dex"

kubectl create namespace dex
kubectl -n dex create secret tls dex.dev.faros.sh.tls --cert=hack/dev/dex/ssl/cert.pem --key=hack/dev/dex/ssl/key.pem

source .env
kubectl -n dex create secret \
    generic github-client \
    --from-literal=client-id=$GITHUB_CLIENT_ID \
    --from-literal=client-secret=$GITHUB_CLIENT_SECRET

kubectl -n dex create secret tls dex-tls --key hack/dev/dex/ssl/key.pem --cert hack/dev/dex/ssl/cert.pem
kubectl apply -f ./hack/dev/dex/dex.yaml
