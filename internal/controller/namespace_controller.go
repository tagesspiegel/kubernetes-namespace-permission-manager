/*
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
*/

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// NamespaceReconciler reconciles a Namespace object
type NamespaceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="*",resources="*",verbs="*"

const (
	LabelNamespacePermissionControl = "ns.tagesspiegel.de/permission-control"
	LabelManagedBy                  = "app.kubernetes.io/managed-by"
	LabelNamespaceName              = "ns.tagesspiegel.de/source-namespace"

	AnnotationNamespaceRoleBindingSubjects = "ns.tagesspiegel.de/rolebinding-subjects"
	AnnotationNamespaceRoleBindingRoleRef  = "ns.tagesspiegel.de/rolebinding-roleref"
	AnnotationNamespaceCustomRoleRules     = "ns.tagesspiegel.de/custom-role-rules"
)

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logx := log.FromContext(ctx)

	ns := &corev1.Namespace{}
	err := r.Client.Get(ctx, req.NamespacedName, ns)
	if err != nil {
		logx.Error(err, "unable to fetch Namespace")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// check if the namespace has our label
	_, ok := ns.Labels[LabelNamespacePermissionControl]
	if !ok {
		logx.V(100).Info("namespace has no label, ignoring")
		// if not, we don't care about it
		return ctrl.Result{}, nil
	}

	roleRef := rbacv1.RoleRef{}

	// check if the namespace has a role ref
	rf, ok := ns.Annotations[AnnotationNamespaceRoleBindingRoleRef]
	if ok {
		rfn, err := ParseRoleBindingRoleRef(rf)
		if err != nil {
			return ctrl.Result{}, err
		}
		roleRef = rfn
	}

	// check if the namespace has custom rules
	cr, ok := ns.Annotations[AnnotationNamespaceCustomRoleRules]
	if ok {
		// parse the role rules
		rules, err := ParseCustomRole(cr)
		if err != nil {
			logx.Error(err, "unable to parse role rules")
			return ctrl.Result{}, nil
		}
		// create a role
		role := &rbacv1.Role{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      ns.Name,
				Namespace: ns.Name,
				Labels: map[string]string{
					LabelManagedBy:     "namespace-permission-controller",
					LabelNamespaceName: ns.Name,
				},
			},
		}
		rslt, err := ctrl.CreateOrUpdate(ctx, r.Client, role, func() error {
			role.Rules = rules
			return nil
		})
		if err != nil {
			logx.Error(err, "unable to create or update role")
			return ctrl.Result{}, nil
		}
		logx.V(80).Info("result for reconciliation for role binding", "result", rslt)
		roleRef.Kind = "Role"
		roleRef.APIGroup = "rbac.authorization.k8s.io"
		roleRef.Name = role.GetName()
	}

	rbSubjects, ok := ns.Annotations[AnnotationNamespaceRoleBindingSubjects]
	if ok {
		// parse the role rules
		subjects, err := ParseRoleBindingSubjects(rbSubjects)
		if err != nil {
			logx.Error(err, "unable to parse role rules")
			return ctrl.Result{}, nil
		}
		// create a rb
		rb := &rbacv1.RoleBinding{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      ns.Name,
				Namespace: ns.Name,
				Labels: map[string]string{
					LabelManagedBy:     "namespace-permission-controller",
					LabelNamespaceName: ns.Name,
				},
			},
		}
		rslt, err := ctrl.CreateOrUpdate(ctx, r.Client, rb, func() error {
			rb.Subjects = subjects
			rb.RoleRef = roleRef
			return nil
		})
		if err != nil {
			logx.Error(err, "unable to create or update rolebinding")
			return ctrl.Result{}, nil
		}
		logx.V(80).Info("result for reconciliation for role binding", "result", rslt)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// we only expect to be called for namespaces with our label
		For(&corev1.Namespace{}, builder.WithPredicates(&LabelChecker{ExpectedLabel: LabelNamespacePermissionControl})).
		Complete(r)
}
