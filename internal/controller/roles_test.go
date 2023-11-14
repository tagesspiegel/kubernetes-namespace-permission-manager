package controller

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	rbacv1 "k8s.io/api/rbac/v1"
)

func TestParseRoleBindingSubjects(t *testing.T) {
	type args struct {
		rulesStr string
	}
	tests := []struct {
		name    string
		args    args
		want    []rbacv1.Subject
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				rulesStr: "kind=ServiceAccount,name=foo,namespace=bar",
			},
			want: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      "foo",
					Namespace: "bar",
				},
			},
			wantErr: false,
		},
		{
			name: "simple with apiGroup",
			args: args{
				rulesStr: "kind=Role,name=foo,namespace=bar,apiGroup=rbac.authorization.k8s.io",
			},
			want: []rbacv1.Subject{
				{
					APIGroup:  "rbac.authorization.k8s.io",
					Kind:      "Role",
					Name:      "foo",
					Namespace: "bar",
				},
			},
			wantErr: false,
		},
		{
			name: "multiple",
			args: args{
				rulesStr: "kind=ServiceAccount,name=foo,namespace=bar;kind=ServiceAccount,name=foo2,namespace=bar2",
			},
			want: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      "foo",
					Namespace: "bar",
				},
				{
					Kind:      "ServiceAccount",
					Name:      "foo2",
					Namespace: "bar2",
				},
			},
			wantErr: false,
		},
		{
			name: "multiple with apiGroup",
			args: args{
				rulesStr: "kind=Role,apiGroup=rbac.authorization.k8s.io,name=foo,namespace=bar;kind=ServiceAccount,name=foo2,namespace=bar2",
			},
			want: []rbacv1.Subject{
				{
					APIGroup:  "rbac.authorization.k8s.io",
					Kind:      "Role",
					Name:      "foo",
					Namespace: "bar",
				},
				{
					Kind:      "ServiceAccount",
					Name:      "foo2",
					Namespace: "bar2",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid key",
			args: args{
				rulesStr: "kind=ServiceAccount,name=foo,namespace=bar;kind=ServiceAccount,name=foo2,namespace=bar2;kind=ServiceAccount,name=foo3,namespace=bar3,invalid=foo",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid amount of entries",
			args: args{
				rulesStr: "kind=ServiceAccount=ServiceAccount,name=foo,namespace=bar;kind=ServiceAccount,name=foo2,namespace=bar2;kind=ServiceAccount,name=foo3,namespace=bar3,invalid=foo",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRoleBindingSubjects(tt.args.rulesStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRoleBindingSubjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("ParseRoleBindingSubjects() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseRoleBindingRoleRef(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    rbacv1.RoleRef
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				str: "kind=Role;apiGroup=rbac.authorization.k8s.io;name=my-role",
			},
			want: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     "my-role",
			},
			wantErr: false,
		},
		{
			name: "simple without apiGroup",
			args: args{
				str: "kind=Role;name=my-role",
			},
			want: rbacv1.RoleRef{
				Kind: "Role",
				Name: "my-role",
			},
			wantErr: false,
		},
		{
			name: "invalid key",
			args: args{
				str: "kind=Role;apiGroup=rbac.authorization.k8s.io;name=my-role;invalid=foo",
			},
			want:    rbacv1.RoleRef{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRoleBindingRoleRef(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRoleBindingRoleRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("ParseRoleBindingRoleRef() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseCustomRole(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    []rbacv1.PolicyRule
		wantErr bool
	}{
		{
			name: "simple success",
			args: args{
				str: "apiGroups=apps;resources=deployments,replicasets;verbs=get,list,watch,create,update,patch,delete",
			},
			want: []rbacv1.PolicyRule{
				{
					APIGroups: []string{"apps"},
					Resources: []string{"deployments", "replicasets"},
					Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
				},
			},
			wantErr: false,
		},
		{
			name: "simple success with resourceNames",
			args: args{
				str: "apiGroups=apps;resources=deployments,replicasets;verbs=get,list,watch,create,update,patch,delete;resourceNames=foo,bar",
			},
			want: []rbacv1.PolicyRule{
				{
					APIGroups:     []string{"apps"},
					Resources:     []string{"deployments", "replicasets"},
					Verbs:         []string{"get", "list", "watch", "create", "update", "patch", "delete"},
					ResourceNames: []string{"foo", "bar"},
				},
			},
			wantErr: false,
		},
		{
			name: "simple success with multiple rules",
			args: args{
				str: "apiGroups=apps;resources=deployments,replicasets;verbs=get,list,watch,create,update,patch,delete::apiGroups=apps;resources=deployments,replicasets;verbs=get,list,watch,create,update,patch,delete",
			},
			want: []rbacv1.PolicyRule{
				{
					APIGroups: []string{"apps"},
					Resources: []string{"deployments", "replicasets"},
					Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
				},
				{
					APIGroups: []string{"apps"},
					Resources: []string{"deployments", "replicasets"},
					Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
				},
			},
			wantErr: false,
		},
		{
			name: "simple success with multiple rules and multiple verbs",
			args: args{
				str: "apiGroups=apps;resources=deployments,replicasets;verbs=get,list,watch,create,update,patch,delete::apiGroups=apps;resources=deployments,replicasets;verbs=get,list,watch,create,update,patch,delete::apiGroups=apps;resources=deployments,replicasets;verbs=get,list,watch,create,update,patch,delete",
			},
			want: []rbacv1.PolicyRule{
				{
					APIGroups: []string{"apps"},
					Resources: []string{"deployments", "replicasets"},
					Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
				},
				{
					APIGroups: []string{"apps"},
					Resources: []string{"deployments", "replicasets"},
					Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
				},
				{
					APIGroups: []string{"apps"},
					Resources: []string{"deployments", "replicasets"},
					Verbs:     []string{"get", "list", "watch", "create", "update", "patch", "delete"},
				},
			},
			wantErr: false,
		},
		{
			name: "simple invalid key value format",
			args: args{
				str: "apiGroups=apps=true;resources=deployments,replicasets;verbs=get,list,watch,create,update,patch,delete",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "simple non existing key",
			args: args{
				str: "apiGroups=apps;resources=deployments,replicasets;verbs=get,list,watch,create,update,patch,delete;invalid=foo",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCustomRole(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCustomRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("ParseCustomRole() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
