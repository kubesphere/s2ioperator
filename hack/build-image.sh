#!/bin/bash
set -e

dest="deploy/s2ioperator.yaml"
tag=`git rev-parse --short HEAD`
IMG=kubespheredev/s2ioperator:$tag

docker build -f deploy/Dockerfile -t ${IMG} bin/
docker push $IMG
echo "updating kustomize image patch file for manager resource"
sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

kustomize build config/default -o $dest
kubectl apply -f $dest
./hack/certs.sh --service webhook-server-service --namespace devops --secret webhook-server-secret