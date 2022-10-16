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

helm repo add nginx-stable https://helm.nginx.com/stable
helm repo update

helm upgrade -i ingress nginx-stable/nginx-ingress \
  --kube-context "${KUBECTL_CONTEXT}" \
  --namespace ingress \
  --create-namespace \
  --set controller.hostPort.enabled=true \
  --set controller.service.type=NodePort \
  --set controller.admissionWebhooks.networkPolicyEnabled=true \
  --set controller.extraArgs.default-ssl-certificate=wayfinder-ingress/ingress-tls-default

echo "Installing cert-manager"

helm repo add jetstack https://charts.jetstack.io
helm repo update

kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.crds.yaml
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.9.1

