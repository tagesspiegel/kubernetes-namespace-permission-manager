#!/bin/bash

set -e

echo "Running e2e tests..."

NAMESPACE_NAME="e2e-test-$(date +%s)"

# existing role in namespace
kubectl create namespace $NAMESPACE_NAME
kubectl annotate namespace $NAMESPACE_NAME ns.tagesspiegel.de/rolebinding-subjects="kind=Group;name=example;apiGroup=rbac.authorization.k8s.io"
kubectl annotate namespace $NAMESPACE_NAME ns.tagesspiegel.de/rolebinding-roleref="apiGroup=rbac.authorization.k8s.io;kind=Role;name=namespace-admin"
kubectl label namespace $NAMESPACE_NAME ns.tagesspiegel.de/permission-control=enabled

wcl=$(kubectl -n $NAMESPACE_NAME get rolebindings.rbac.authorization.k8s.io -l app.kubernetes.io/managed-by=namespace-permission-controller -l ns.tagesspiegel.de/source-namespace=$NAMESPACE_NAME -o json | jq -r ".items[].metadata.name" out.json | wc -l)

if [ $wcl -ne 1 ]; then
  echo "Expected 1 rolebinding, got $wcl"
  exit 1
fi
