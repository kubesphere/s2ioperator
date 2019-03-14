#!/bin/bash
set -e

function cleanup(){
    result=$?
    echo "Cleaning"
    kubectl delete ns $TEST_NS
    exit $result
}
dest="./deploy/s2ioperator.yaml"
tag=`git rev-parse --short HEAD`
IMG=kubespheredev/s2ioperator:$tag
TEST_NS=s2ioperator-test-$tag

trap cleanup EXIT SIGINT SIGQUIT
docker build -f deploy/Dockerfile -t ${IMG} bin/
docker push $IMG
echo "updating kustomize image patch file for manager resource"

if [ "$(uname)" == "Darwin" ]; then
    sed -i '' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml
    sed -i '' -e  's/namespace: .*/namespace: '"${TEST_NS}"'/' ./config/default/kustomization.yaml
else
    sed -i  -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml
    sed -i  -e  's/namespace: .*/namespace: '"${TEST_NS}"'/' ./config/default/kustomization.yaml
fi
kubectl create ns  $TEST_NS
kustomize build config/default -o $dest
kubectl apply -f $dest
./hack/certs.sh --service webhook-server-service --namespace $TEST_NS --secret webhook-server-secret

export TEST_NS
go test -v ./test/e2e/
