
# Image URL to use all building/pushing image targets
IMG ?= kubespheredev/s2ioperator:latest
NAMESPACE ?= kubesphere-devops-system
export GO111MODULE=on

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

all: test manager

# Run tests
test: fmt vet
	export KUBEBUILDER_CONTROLPLANE_START_TIMEOUT=1m; ginkgo -v -cover ./pkg/...  

# Build manager binary
manager: generate fmt manifests vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -o bin/manager github.com/kubesphere/s2ioperator/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/manager/main.go

# Install CRDs into a cluster
install-crd: manifests
	kubectl apply -f config/crds

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests update-cert kustomize
	kubectl apply -f config/templates
	kubectl apply -f config/crds
	$(KUSTOMIZE) config | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) crd:trivialVersions=true rbac:roleName=manager-role paths="./pkg/apis/...;./pkg/controller/..." output:crd:artifacts:config=config/crds

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

client-gen:
	./hack/generate_group.sh all github.com/kubesphere/s2ioperator/pkg/client github.com/kubesphere/s2ioperator/pkg/apis "devops:v1alpha1" --go-header-file ./hack/boilerplate.go.txt

# Generate code
generate: controller-gen openapi-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./pkg/apis/...
	$(OPENAPI_GEN) -O openapi_generated -i k8s.io/api/core/v1,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/api/resource,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/util/intstr,k8s.io/apimachinery/pkg/version,github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1 -p github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1 -h hack/boilerplate.go.txt --report-filename api/api-rules/violation_exceptions.list

# Build the docker image
docker-build:
	docker build -f deploy/Dockerfile -t $(IMG) bin/
	docker push $(IMG)

image-multiarch:
	docker build . -t $(IMG) -f deploy/Dockerfile.multiarch

debug: manager
	./hack/build-image.sh

release: kustomize manager test docker-build update-cert
	$(KUSTOMIZE) config > deploy/s2ioperator.yaml

install-travis:
	chmod +x ./hack/*.sh
	./hack/install_tools.sh

e2e-test:
	./hack/e2etest.sh

# create the secret with CA cert and server cert/key
ca-secret:
	./hack/certs.sh --service webhook-server-service --namespace $(NAMESPACE)

# update certs
update-cert: ca-secret
	./hack/update-cert.sh

.PHONY : clean test

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.6.2)

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

OPENAPI_GEN = $(shell pwd)/bin/openapi-gen
openapi-gen: ## Download openapi-gen locally if necessary
	$(call go-get-tool,$(OPENAPI_GEN),k8s.io/code-generator/cmd/openapi-gen@v0.22.3)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
