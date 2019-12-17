#!/bin/bash

set -e

usage() {
    cat <<EOF
Generate certificate suitable for use with an sidecar-injector webhook service.
This script uses k8s' CertificateSigningRequest API to a generate a
certificate signed by k8s CA suitable for use with sidecar-injector webhook
services. This requires permissions to create and approve CSR. See
https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster for
detailed explantion and additional instructions.
The server key/cert k8s CA cert are stored in a k8s secret.
usage: ${0} [OPTIONS]
The following flags are required.
       --service          Service name of webhook.
       --namespace        Namespace where webhook service and secret reside.
       --secret           Secret name for CA certificate and server certificate/key pair.
EOF
    exit 1
}

while [[ $# -gt 0 ]]; do
    case ${1} in
        --service)
            service="$2"
            shift
            ;;
        --secret)
            secret="$2"
            shift
            ;;
        --namespace)
            namespace="$2"
            shift
            ;;
        *)
            usage
            ;;
    esac
    shift
done

[ -z ${service} ] && service=sidecar-injector-webhook-svc
[ -z ${secret} ] && secret=sidecar-injector-webhook-certs
[ -z ${namespace} ] && namespace=default

if [ ! -x "$(command -v openssl)" ]; then
    echo "openssl not found"
    exit 1
fi

csrName=${service}.${namespace}
certsdir="config/certs"

if [ ! -d ${certsdir} ]; then
  mkdir ${certsdir}
fi

echo "creating certs in certsdir ${certsdir} "

cat <<EOF >> ${certsdir}/csr.conf
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = ${service}
DNS.2 = ${service}.${namespace}
DNS.3 = ${service}.${namespace}.svc
EOF

openssl genrsa -out ${certsdir}/server-key.pem 2048
openssl req -new -key ${certsdir}/server-key.pem -subj "/CN=${service}.${namespace}.svc" -out ${certsdir}/server.csr -config ${certsdir}/csr.conf

# clean-up any previously created CSR for our service. Ignore errors if not present.
kubectl delete csr ${csrName} 2>/dev/null || true

# create  server cert/key CSR and  send to k8s API
cat <<EOF | kubectl create -f -
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: ${csrName}
spec:
  groups:
  - system:authenticated
  request: $(cat ${certsdir}/server.csr | base64 | tr -d '\n')
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF

# verify CSR has been created
while true; do
    kubectl get csr ${csrName}
    if [ "$?" -eq 0 ]; then
        break
    fi
done

# approve and fetch the signed certificate
kubectl certificate approve ${csrName}
# verify certificate has been signed
for x in $(seq 10); do
    serverCert=$(kubectl get csr ${csrName} -o jsonpath='{.status.certificate}')
    if [[ ${serverCert} != '' ]]; then
        break
    fi
    sleep 1
done
if [[ ${serverCert} == '' ]]; then
    echo "ERROR: After approving csr ${csrName}, the signed certificate did not appear on the resource. Giving up after 10 attempts." >&2
    exit 1
fi
echo ${serverCert} | openssl base64 -d -A -out ${certsdir}/server-cert.pem

kubectl config view --raw -o json | jq -r '.clusters[0].cluster."certificate-authority-data"' | tr -d '"' | base64 --decode > ${certsdir}/ca.pem
# create the secret with CA cert and server cert/key
kubectl create secret generic ${secret} \
        --from-file=tls.key=${certsdir}/server-key.pem \
        --from-file=tls.crt=${certsdir}/server-cert.pem \
        --from-file=ca.crt=${certsdir}/ca.pem \
        --dry-run -o yaml |
    kubectl -n ${namespace} apply -f -

muWebhook=$(kubectl get mutatingwebhookconfigurations mutating-webhook-configuration -o )
if [[ ${muWebhook} != '' ]]; then
    cabundle=$(cat ${certsdir}/ca.pem | base64)

    kubectl patch mutatingwebhookconfigurations mutating-webhook-configuration --type='json' -p="[{\"op\": \"replace\", \"path\": \"/webhooks/0/clientConfig/caBundle\", \"value\":\"${cabundle}\"}]"

    for ((i=0;i<=2;i++)); do
        kubectl patch validatingwebhookconfigurations validating-webhook-configuration --type='json' -p="[{\"op\": \"replace\", \"path\": \"/webhooks/${i}/clientConfig/caBundle\", \"value\":\"${cabundle}\"}]"
    done
    fi