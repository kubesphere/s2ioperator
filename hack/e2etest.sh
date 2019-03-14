#!/bin/bash
set -e

function cleanup(){
    result=$?
    echo "Cleaning"
    kubectl delete ns $TEST_NS
    exit $result
}
dest="/tmp/s2ioperator.yaml"
tag=`git rev-parse --short HEAD`
IMG=kubespheredev/s2ioperator:$tag
TEST_NS=s2ioperator-test-$tag

trap cleanup EXIT SIGINT SIGQUIT
docker build -f deploy/Dockerfile -t ${IMG} bin/
docker push $IMG
echo "updating kustomize image patch file for manager resource"
sed -i '' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

kubectl create ns  $TEST_NS
sed -i '' -e  's/namespace: .*/namespace: '"${TEST_NS}"'/' ./config/default/kustomization.yaml
kustomize build config/default -o $dest
kubectl apply -f $dest
./hack/certs.sh --service webhook-server-service --namespace $TEST_NS --secret webhook-server-secret

export TEST_NS
go test -v ./test/e2e/
