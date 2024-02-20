#!/bin/bash

set -e

CaBundle=$(< ./config/certs/ca.crt base64 -w 0)
TLSKey=$(< ./config/certs/server.key base64 -w 0)
TLSCrt=$(< ./config/certs/server.crt base64 -w 0)

echo "Update Secret: s2i-webhook-server-cert.."
kubectl -n kubesphere-devops-system patch secret s2i-webhook-server-cert --type='json' -p="[\
{\"op\": \"replace\", \"path\": \"/data/caBundle\", \"value\": \"${CaBundle}\"},\
{\"op\": \"replace\", \"path\": \"/data/tls.key\", \"value\": \"${TLSKey}\"},\
{\"op\": \"replace\", \"path\": \"/data/tls.crt\", \"value\": \"${TLSCrt}\"}\
]"

echo "Update ValidatingWebhookConfiguration validating-webhook-configuration.."
kubectl -n kubesphere-devops-system patch validatingwebhookconfigurations validating-webhook-configuration --type='json' -p="[\
{\"op\": \"replace\", \"path\": \"/webhooks/0/clientConfig/caBundle\", \"value\": \"${CaBundle}\"},\
{\"op\": \"replace\", \"path\": \"/webhooks/1/clientConfig/caBundle\", \"value\": \"${CaBundle}\"},\
{\"op\": \"replace\", \"path\": \"/webhooks/2/clientConfig/caBundle\", \"value\": \"${CaBundle}\"}\
]"

echo "Update MutatingWebhookConfiguration mutating-webhook-configuration.."
kubectl -n kubesphere-devops-system patch mutatingwebhookconfigurations mutating-webhook-configuration --type='json' -p="[{\"op\": \"replace\", \"path\": \"/webhooks/0/clientConfig/caBundle\", \"value\": \"${CaBundle}\"}]"

echo "Restart s2ioperator server.."
sleep 5
kubectl -n kubesphere-devops-system rollout restart sts s2ioperator


echo "Done."
