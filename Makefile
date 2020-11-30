
# Image URL to use all building/pushing image targets
IMG ?= kubespheredev/s2ioperator:latest
NAMESPACE ?= kubesphere-devops-system
export GO111MODULE=on

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
deploy: manifests update-cert
	kubectl apply -f config/crds
	kubectl kustomize config | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go  crd:trivialVersions=true rbac:roleName=manager-role paths="./pkg/apis/...;./pkg/controller/..." output:crd:artifacts:config=config/crds

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

client-gen:
	./hack/generate_group.sh all github.com/kubesphere/s2ioperator/pkg/client github.com/kubesphere/s2ioperator/pkg/apis "devops:v1alpha1" --go-header-file ./hack/boilerplate.go.txt

# Generate code
generate:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go  object:headerFile=./hack/boilerplate.go.txt paths=./pkg/apis/...
	go run vendor/k8s.io/kube-openapi/cmd/openapi-gen/openapi-gen.go -O openapi_generated -i k8s.io/api/core/v1,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/api/resource,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/util/intstr,k8s.io/apimachinery/pkg/version,github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1 -p github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1 -h hack/boilerplate.go.txt --report-filename api/api-rules/violation_exceptions.list


# Build the docker image
docker-build:
	docker build -f deploy/Dockerfile -t $(IMG) bin/
	docker push $(IMG)

debug: manager
	./hack/build-image.sh

release: manager test docker-build update-cert
	kubectl kustomize config > deploy/s2ioperator.yaml

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
