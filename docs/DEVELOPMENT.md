# S2I Operator Development Guide

## Getting started

1. [Ramp up on kubernetes and CRDs](DEVELOPMENT.md#ramp-up-on-kubernetes-and-crds)
2. [Checkout your fork](DEVELOPMENT.md#checkout-your-fork)
3. [Set up your develop environment](DEVELOPMENT.md#set-up-local-develop-environment)
4. [Iterating](DEVELOPMENT.md#iterating)
5. [Test](DEVELOPMENT.md#test)
6. [Install CRD](DEVELOPMENT.md#install-s2i-crd)

### Ramp up on kubernetes and CRDs

Welcome to the project. S2I Operator are work on Kubernetes, and all steps about ci/cd are defined by [CustomResourceDefinitions](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/) .Here are some resources helpful to ramp up on some of the technology this project is built on.

- [Understanding Kubernetes objects](https://kubernetes.io/docs/concepts/overview/working-with-objects/kubernetes-objects/) 
- [API conventions - Types(kinds)](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#types-kinds) 
- [Extend the Kubernetes API with CustomResourceDefinitions](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/)

### Checkout your fork

Checkout this repository:

1. Create your own fork  of [this repo](https://github.com/kubesphere/s2ioperator).

2. Change to your work directory, and clone it

   ```shell
   git clone https://github.com/${YOUR_GITHUB_USERNAME}/s2ioperator.git
   cd s2ioperator
   git remote add upstream git@github:kubesphere/s2ioperator.git
   git remote set-url --push upstream no_push
   ```

### Set up local develop environment

#### Prerequisites

- [golang](https://golang.org/dl/) environment
- [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) 2.0+.
- [docker](https://docs.docker.com/install/) version 17.03+.
- [kustomize](https://sigs.k8s.io/kustomize/docs/INSTALL.md) v3.1.0+

Also you can install some tools by script `install_tools.sh` in `hack`:

```shell
./hack/install_tools.sh
```

#### Build S2I Operator locally

1. Set environment variable `GO111MODULE=auto` , if directory of s2ioperator in your [`GOPATH`](https://github.com/golang/go/wiki/SettingGOPATH), set `GO111MODULE=on`

2. Run the following command to create a binary with the source code:

   ```shell
   make manager
   ```

If nothing goes wrong, and output a binary in directory `WORKDIR/bin/`, mean local environment is ok.

#### Configure kubectl to use your cluster

To debug your S2I Operator in local, you will need to [Configure kubectl to use your kubernetes cluster](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/).

If you don't have a cluster, please reference [Kubernetes doc](https://kubernetes.io/docs/setup/).

### Iterating

While iterating on the project, you may need to:

1. Running [unit test and end-to-end](DEVELOPMENT.md#test) to ensure your code are works well.

2. Update your (external) dependencies with: `./hack/update-vendor.sh`

   Also you should running `go mod verify` to verify dependencies have expected content.

3. Running `make deploy` to deploy S2I Operator and verify it's working by looking at the pod logs.

4. Running `make release` to release your code, and commit by git.

To make changes to these CRDs, you will probably interact with

- The CRD type definitions in [./pkg/apis/devops/alpha1](https://github.com/kubesphere/s2ioperator/tree/master/pkg/apis/devops/v1alpha1)
- The reconcilers in [./pkg/controller](https://github.com/kubesphere/s2ioperator/tree/master/pkg/controller)
- The clients are in [./pkg/client](https://github.com/kubesphere/s2ioperator/tree/master/pkg/client) 

### Test

Before run test, you shoud install test tools `ginkgo` by following command:

```go
go get -u  github.com/onsi/ginkgo/ginkgo
```

#### Unit tests

Unit tests live side by side with the code they are testing and will run `go fmt`and `go vet` before test by default. You can run unit test with:

```shell
make test
```

#### End to end tests

By default the tests run against your current kubeconfig context. You  can run e2e test with:

```shell
make e2e-test
```

**Node:  Running End to end tests will change default configs in the directory `config`**

### Install S2I CRD

All CustomResourceDefinitions(CRD) used in S2I Operator are defined in `config/crds/`. You can use command `make manifests` to generate manifests, e.g. CRD, RBAC etc.

Install S2I CRD with following command:

```shell
make install-crd
```

For more information about S2I CRD please to see [here](CRD-Consepts.md).
