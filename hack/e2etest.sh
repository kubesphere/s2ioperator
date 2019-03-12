#!/bin/bash
set -e

dest="deploy/s2ioperator.yaml"
tag=`git rev-parse --short HEAD`
IMG=kubespheredev/s2ioperator:$tag
TEST_NS=s2ioperator-test-$tag

docker build -f deploy/Dockerfile -t ${IMG} bin/
docker push $IMG
echo "updating kustomize image patch file for manager resource"
sed -i -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

kubectl create ns  $TEST_NS
sed -i -e 's/namespace: .*/namespace: '"${TEST_NS}"'/' ./config/default/kustomization.yaml
kustomize build config/default -o $dest
kubectl apply -f $dest
./hack/certs.sh --service webhook-server-service --namespace $TEST_NS --secret webhook-server-secret

set +e
export TEST_NS
go test -v ./test/e2e/
result=$?
kubectl delete ns $TEST_NS
exit $result