REPO ?= quay.io/faroshq/kcp-potatoes-service
TAG_NAME ?= $(shell git describe --tags --abbrev=0)
LOCALBIN ?= $(shell pwd)/bin
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
KUSTOMIZE ?= $(LOCALBIN)/kustomize

KUSTOMIZE_VERSION ?= v3.8.7
CONTROLLER_TOOLS_VERSION ?= v0.8.0

# KCP prefix
APIEXPORT_PREFIX ?= v$(shell date +'%Y%m%d')

run-server:
	@echo "Starting server..."
	@go run ./cmd/server

run-server-mod:
	@echo "Starting server..."
	@go run -mod=vendor ./cmd/server

build-server:
	@echo "Building server..."
	@go build -o ./bin/server ./cmd/server

build-potatoes:
	@echo "Building potatoes..."
	@go build -o ./bin/potatoes ./cmd/potatoes

.PHONY: image-server
image-server:
	docker build -t ${REPO}:${TAG_NAME} -f dockerfiles/server/Dockerfile .

.PHONY: image-potatoes
image-potatoes:
	docker build -t ${REPO}:${TAG_NAME} -f dockerfiles/potatoes/Dockerfile .

images: image-server image-potatoes

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
$(KUSTOMIZE): ## Download kustomize locally if necessary.
	mkdir -p $(LOCALBIN)
	curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN)
	touch $(KUSTOMIZE) # we download an "old" file, so make will re-download to refresh it unless we make it newer than the owning dir

$(CONTROLLER_GEN): ## Download controller-gen locally if necessary.
	mkdir -p $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

manifests: $(CONTROLLER_GEN) ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./pkg/controllers/..." output:crd:artifacts:config=config/crd/bases

.PHONY: apiresourceschemas
apiresourceschemas: $(KUSTOMIZE) ## Convert CRDs from config/crds to APIResourceSchemas. Specify APIEXPORT_PREFIX as needed.
	$(KUSTOMIZE) build config/crd | kubectl kcp crd snapshot -f - --prefix $(APIEXPORT_PREFIX) > config/kcp/$(APIEXPORT_PREFIX).apiresourceschemas.yaml


.PHONY: generate
generate: $(CONTROLLER_GEN) ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./pkg/controllers/..."
