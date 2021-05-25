[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/kubesphere/s2ioperator)

# Source-to-image Operator

[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/kubesphere/s2ioperator/blob/master/LICENSE)  [![Go Report Card](https://goreportcard.com/badge/github.com/kubesphere/s2ioperator)](https://goreportcard.com/report/github.com/kubesphere/s2ioperator)  [![S2I Operator release](https://img.shields.io/github/release/kubesphere/s2ioperator.svg?color=release&label=release&logo=release&logoColor=release)](https://github.com/kubesphere/s2ioperator/releases/tag/v0.0.14)

Source-to-image(S2I)-Operator is a Kubernetes [Custom Resource Defintion](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) (CRD) controller that provides easy Kubernetes-style resources for declaring CI/CD-style pipelines. S2I Operator create a ready-to-run images by injecting source code into a container image and letting the container prepare that source code for execution. And create workload automatically with ready-to-run images.

## Native S2I vs S2I Operator

Compare with native S2I, S2I Operator also has the following advantages on the original foundation.

1. **Provide S2I Open API**: you can call S2I API directly to generate image, or carry out secondary development.
2. **Simple Config**: you just provide source code url, and specify the image repo which  you want to push, other configurations will setted automatically. And all configurations are stored as different resources in Kubernetes.
3. **Deep integration with Kubernetes**: Use containers as their building blocks. And you can use kubectl to create s2i pipelines just as you do with Kubernetes' built-in resources.

## Installation

#### Prerequisites

1. A Kubernetes cluster. (if you don't have an existing cluster, please [create it](https://kubernetes.io/docs/setup/).
2. Grant cluster-admin permissions to the current user.

#### Install S2I Operator

You can install S2I Operator in any kubernetes cluster with following commands:

```shell
# create a namespaces, such as kubesphere-devops-system
kubectl create ns kubesphere-devops-system
# create S2I Operator and all CRD 
kubectl apply -f  https://github.com/kubesphere/s2ioperator/releases/download/v0.0.2/s2ioperator.yaml
```

Now monitor the S2I Operator components show a `STATUS` of `Running`:

```shell
# please change you namespace
kubectl -n kubesphere-devops-system get pods -w
```

## Quick Start

Here is [quick-start](docs/QUICK-START.md) to walk you through the process, with a quick overview of the core features of S2I Operator that helps you to get familiar with it.

If you want to get a better experience with S2I Operator, perhaps you can use S2I CI/CD in [Kubesphere](https://github.com/kubesphere/kubesphere).

## Welcome to contribute

We are so excited to have you!

- See [Kubesphere community guide](https://github.com/kubesphere/community) for an overview of our processes
- See [DEVELOPMENT.md](docs/DEVELOPMENT.md) for how to get started
