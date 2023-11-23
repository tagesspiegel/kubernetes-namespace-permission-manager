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

## Getting Started

### Prerequisites

- Go version v1.21.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/kubernetes-namespace-permission-manager:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/kubernetes-namespace-permission-manager:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall

**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
