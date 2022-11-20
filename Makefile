REPO ?= quay.io/faroshq/
TAG_NAME ?= $(shell git describe --tags --abbrev=0)
LOCALBIN ?= $(shell pwd)/bin
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
GO_INSTALL = ./hack/go-install.sh
KUSTOMIZE ?= $(LOCALBIN)/kustomize
TOOLS_DIR=hack/tools
TOOLS_GOBIN_DIR := $(abspath $(TOOLS_DIR))
KO_DOCKER_REPO ?= ${REPO}

KUSTOMIZE_VERSION ?= v3.8.7
CONTROLLER_GEN_VER := v0.10.0
CONTROLLER_GEN_BIN := controller-gen

CONTROLLER_GEN := $(TOOLS_DIR)/$(CONTROLLER_GEN_BIN)-$(CONTROLLER_GEN_VER)
export CONTROLLER_GEN # so hack scripts can use it

# KCP prefix
#APIEXPORT_PREFIX ?= v$(shell date +'%Y%m%d')
APIEXPORT_PREFIX = today

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
$(KUSTOMIZE): ## Download kustomize locally if necessary.
	mkdir -p $(LOCALBIN)
	curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN)
	touch $(KUSTOMIZE) # we download an "old" file, so make will re-download to refresh it unless we make it newer than the owning dir

manifests:  ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crds/bases | true
	make generate

.PHONY: apiresourceschemas
apiresourceschemas: $(KUSTOMIZE) ## Convert CRDs from config/crds to APIResourceSchemas. Specify APIEXPORT_PREFIX as needed.
	$(KUSTOMIZE) build config/crds | kubectl kcp crd snapshot -f - --prefix $(APIEXPORT_PREFIX) > config/kcp/$(APIEXPORT_PREFIX).apiresourceschemas.yaml
	make generate

tools:$(CONTROLLER_GEN)
.PHONY: tools

$(CONTROLLER_GEN):
	GOBIN=$(TOOLS_GOBIN_DIR) $(GO_INSTALL) sigs.k8s.io/controller-tools/cmd/controller-gen $(CONTROLLER_GEN_BIN) $(CONTROLLER_GEN_VER)

codegen: $(CONTROLLER_GEN) generate ## Run the codegenerators
	go mod download
	./hack/update-codegen.sh
.PHONY: codegen

generate:
	go generate ./...

lint:
	gofmt -s -w cmd hack pkg
	go run golang.org/x/tools/cmd/goimports -w -local=github.com/faroshq/faros-hub cmd hack pkg
	go run ./hack/validate-imports cmd hack pkg
	staticcheck ./...

setup-kind:
	./hack/dev/setup-kind.sh

delete-kind:
	./hack/dev/delete-kind.sh

run-with-oidc:
	./hack/dev/run-with-oidc.sh

images:
	KO_DOCKER_REPO=${KO_DOCKER_REPO} ko build --sbom=none -B --platform=linux/amd64 -t latest ./cmd/*
