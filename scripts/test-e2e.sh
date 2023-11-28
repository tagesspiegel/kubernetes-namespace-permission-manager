#!/bin/bash

set -e

echo "Running e2e tests..."

color_red="\e[31m"
color_green="\e[32m"
color_blue="\e[34m"
color_reset="\e[0m"
color_yellow="\e[33m"

# deleteNamespace <namespace>
function deleteNamespace() {
  kubectl delete namespace $1
}

# generateName
function generateName() {
  echo "e2e-test-$(date +%s%3N)"
}

# createNamespace <namespace> <rolebinding-subjects>
function createNamespace() {
  kubectl create namespace $1
  kubectl annotate namespace $1 ns.tagesspiegel.de/rolebinding-subjects="$2"
  kubectl label namespace $1 ns.tagesspiegel.de/permission-control=enabled
}

# withRoleRef <namespace> <role>
function withRoleRef() {
  kubectl annotate namespace $1 ns.tagesspiegel.de/rolebinding-roleref="$2"
}

# withCustomRules <namespace> <custom-role-rules>
function withCustomRules() {
  kubectl annotate namespace $1 ns.tagesspiegel.de/custom-role-rules="$2"
}

# createRole <namespace> <role> "<verb>" "<resource>"
function createRole() {
  kubectl -n $1 create role $2 --verb=$3 --resource=$4
}

# checkForRoleBinding <namespace> <rolebinding>
function checkForRoleBinding() {
  local name=$(kubectl -n $1 get rolebindings.rbac.authorization.k8s.io $1 -o json | jq '.metadata.name')
  if [ "$name" != "\"$2\"" ]; then
    echo 1
  fi
  echo 0
}

# assert <namespace> <got> <expected>
# check if got is equal to expected, if not, delete namespace and exit
function assert() {
  if [ "$2" != "$3" ]; then
    echo -e "${color_red}Test for ${color_blue}$1 ${color_red}failed!${color_red}\n\tExpected ${color_yellow}$3${color_red}, got ${color_yellow}$2${color_reset}"
    deleteNamespace $1
    exit 1
  fi
  #deleteNamespace $1
  echo -e "${color_green}Test for ${color_yellow}$1${color_green} passed successfully!${color_reset}"
}

echo "========== Test 1 =========="
NAMESPACE_NAME=$(generateName)
ROLE_NAME="namespace-admin"
createNamespace $NAMESPACE_NAME "kind=Group;name=example"
createRole $NAMESPACE_NAME $ROLE_NAME "get,watch,list" "pods"
withRoleRef $NAMESPACE_NAME "kind=Role;apiGroup=rbac.authorization.k8s.io;name=$ROLE_NAME"
# the operator needs some time to pick up the namespace change
sleep 10
rslt=$(checkForRoleBinding $NAMESPACE_NAME $NAMESPACE_NAME)
assert $NAMESPACE_NAME $rslt "0"

echo "========== Test 2 =========="
NAMESPACE_NAME=$(generateName)
createNamespace $NAMESPACE_NAME "kind=Group;name=example"
withCustomRules $NAMESPACE_NAME "verbs=*;apiGroups=*;resources=*"
# the operator needs some time to pick up the namespace change
sleep 10
rslt=$(checkForRoleBinding $NAMESPACE_NAME $NAMESPACE_NAME)
assert "$NAMESPACE_NAME" "$rslt" "0"
