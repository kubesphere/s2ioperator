
# Image URL to use all building/pushing image targets
IMG ?= kubespheredev/s2ioperator:v0.0.1

all: test manager

# Run tests
test: fmt vet
	ginkgo -v -cover ./pkg/...  

# Build manager binary
manager: generate fmt manifests vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -o bin/manager github.com/kubesphere/s2ioperator/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/manager/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crds

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crds
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate:
	go run vendor/k8s.io/code-generator/cmd/deepcopy-gen/main.go -O zz_generated.deepcopy -i github.com/kubesphere/s2ioperator/pkg/apis/... -h hack/boilerplate.go.txt

# Build the docker image
docker-build: 
	docker build -f deploy/Dockerfile -t $(IMG) bin/
	docker push $(IMG)
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

debug: manager
	./hack/build-image.sh
release: manager test docker-build
	kustomize build config/default -o deploy/s2ioperator.yaml

install-travis:
	chmod +x ./hack/*.sh
	./hack/install_tools.sh

e2e-test: debug
	ginkgo -v  ./e2e/