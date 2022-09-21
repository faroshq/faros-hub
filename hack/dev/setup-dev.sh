#!/bin/bash


if [ ! -f "/usr/local/bin/kind" ]; then
 echo "Installing KIND"
 curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.14.0/kind-linux-amd64
 chmod +x ./kind
 sudo mv ./kind /usr/local/bin/kind
else
    echo "KIND already installed"
fi

kind create cluster --name cluster1 --kubeconfig ./dev/cluster1
kind create cluster --name cluster2 --kubeconfig ./dev/cluster2

echo "Configure faros dev"

kubectl faros workload sync kind1 --syncer-image ghcr.io/kcp-dev/kcp/syncer:main -o ./dev/syncer-kind1-main.yaml
kubectl faros workload sync kind2 --syncer-image ghcr.io/kcp-dev/kcp/syncer:main -o ./dev/syncer-kind2-main.yaml

KUBECONFIG=./dev/cluster1 kubectl apply -f "./dev/syncer-kind1-main.yaml"
KUBECONFIG=./dev/cluster2 kubectl apply -f "./dev/syncer-kind2-main.yaml"
