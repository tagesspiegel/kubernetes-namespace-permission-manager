# kubernetes-namespace-permission-manager

The kubernetes-namespace-permission-manager is a Kubernetes controller that manages the permissions of the namespaces in the cluster. We at Tagesspiegel use ArgoCD to generate feature environments based on application-sets. Since we create a dedicated namespace per environment, we need to manage the permissions to the namespaces. This controller is responsible for managing the permissions to the namespaces, based on annotations passed to the namespace.

## Description

In order to manage the permissions to the namespaces, you need to decide if you want to create a "custom-role" in that namespace or use an pre-existing one. Depending on that, you need to pass the following annotations to the namespace:

| Annotation | Description |
|---|---|
| `ns.tagesspiegel.de/rolebinding-subjects` | A comma separated list of subjects that should be bound to the role. We expect key=value pairs in every array index seperated by semicolons. Example: `a=b;c=d,a=c;b=d`. Valid property keys are: `kind`, `name`, `namespace`. |
| `ns.tagesspiegel.de/rolebinding-roleref` | Semicolon seperated key=value pairs. Example: `a=b;c=d,a=c;b=d`. Valid property keys are: `kind`, `apiGroup`, `name`. |
| `ns.tagesspiegel.de/custom-role-rules` | A two colon `::` seperated list of policy properties, attached to the custom Role. Every array entry is expected to have the following key=value specifications: </br>key=`verbs` a comma seperated list of policy verbs (like: `get`, `list`, `watch`, `patch`, `update`, `delete`, `create`, ...)</br>key=`apiGroups` as list of comma seperated apis to grant access to</br>key=`resources` a list of comma seperated api resources to grant access to</br>key=`resourceNames` (optional) as list of comma seperated resources to grant access to.</br></br>Has priority over `ns.tagesspiegel.de/rolebinding-roleref` |

Since these annotations are not in charge of instrumenting the controller to listen to the namespace, you need to add the following label to the namespace:

| Label | Description |
|---|---|
| `ns.tagesspiegel.de/permission-control` | The value of this label is not important. It is just used to identify the namespaces that should be managed by the controller. |

## Installation

### Using Helm

We recommend using [Helm](https://helm.sh/) to install the controller. We at Tagesspiegel provide a versioned Helm Chart for this controller. You can find the Helm Chart [here](https://github.com/tagesspiegel/helm-charts/tree/main/charts/namespace-permission-manager). In order to install the controller using Helm, you can run the following command:

```bash
helm repo add tagesspiegel https://tagesspiegel.github.io/helm-charts
helm upgrade --install my-name tagesspiegel/namespace-permission-manager
```

For more information about the Helm Chart, please visit the [Helm Chart repository](https://github.com/tagesspiegel/helm-charts/tree/main/charts/namespace-permission-manager).

### Kustomize

Since this repository is using [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) to bootstrap the controller, you can also use [Kustomize](https://kustomize.io/) to install the controller. Please keep in mind that the Kustomize files might include breaking changes, since they are not part of the release process.

In order to install the controller using Kustomize, you can run the following command:

```bash
export IMG=tagesspiegel/kubernetes-namespace-permission-manager:<version>
make deploy
```

## Getting Started

### Prerequisites

- Go version v1.21.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### Testing

In order to test the controller locally, you need to have a Kubernetes cluster running. You can use [kind](https://kind.sigs.k8s.io/) to create a local cluster. Once you have a cluster running, you can run the following command to run the controller locally:

```bash
make run
```

#### Creating namespaces with the required annotations

In order to test the controller, you need to create namespaces with the required annotations and label. You can use the following command to create two namespaces with the available annotations and label:

```bash
kubectl apply -f config/samples/
```

This will create two namespaces with the names `with-roleref` and `with-custom-role` and the required annotations and label.

## Contributing

We welcome contributions. Please follow the [contribution guide](CONTRIBUTING.md) to get started. If you contribute to this project, you agree to abide by its [code of conduct](CODE_OF_CONDUCT.md).
