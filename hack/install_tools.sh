#!/bin/bash

set -e

echo "install kubebuilder"


version=2.3.0 # edit me.

os=$(go env GOOS)
arch=$(go env GOARCH)

# download the release
curl -L -O "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${version}/kubebuilder_${version}_${os}_${arch}.tar.gz"

# extract the archive
tar -zxvf kubebuilder_${version}_${os}_${arch}.tar.gz
sudo mv kubebuilder_${version}_${os}_${arch} /usr/local/kubebuilder

# update your PATH to include /usr/local/kubebuilder/bin
export PATH=$PATH:/usr/local/kubebuilder/bin

# echo "install kustomize"

# wget https://github.com/kubernetes-sigs/kustomize/releases/download/v1.0.11/kustomize_1.0.11_linux_amd64 
# chmod u+x kustomize_1.0.11_linux_amd64
# mv kustomize_1.0.11_linux_amd64 /home/travis/bin/kustomize
echo "install test tools"
go get -u  github.com/onsi/ginkgo/ginkgo

echo "Tools install done"