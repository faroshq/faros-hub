#!/bin/bash


if [ ! -f "/usr/local/bin/kind" ]; then
 echo "Installing KIND"
 curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.14.0/kind-linux-amd64
 chmod +x ./kind
 sudo mv ./kind /usr/local/bin/kind
else
    echo "KIND already installed"
fi

kind create cluster --name services1 --kubeconfig ./dev/services1.kubeconfig
kind create cluster --name shared1 --kubeconfig ./dev/shared1.kubeconfig
