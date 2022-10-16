#!/bin/bash

export ARCH=$(go env GOARCH)
export OS=$(go env GOOS)

FAROS_API_IMAGE=$(KO_DOCKER_REPO=kind.local KIND_CLUSTER_NAME=faros ko build --platform=linux/$ARCH ./cmd/hub-api)
export KUBECONFIG=$(pwd)/dev/faros.kubeconfig
cd manifest && kustomize edit set image faros=$FAROS_API_IMAGE && kustomize build . | kubectl apply -f -
