package controller

import (
	"errors"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"

	strc "github.com/tagesspiegel/kubernetes-namespace-permission-manager/utils/strings"
)

var (
	ErrInvalidKeyInRole       = errors.New("invalid key in role")
	ErrInvalidKeyInRoleRef    = errors.New("invalid key in role ref")
	ErrInvalidKeyInCustomRole = errors.New("invalid key in custom role")
)

const (
	KeyKind          = "kind"
	KeyAPIGroup      = "apiGroup"
	KeyName          = "name"
	KeyNamespace     = "namespace"
	KeyVerbs         = "verbs"
	KeyAPIGroups     = "apiGroups"
	KeyResources     = "resources"
	KeyResourceNames = "resourceNames"
)

// ParseRoleBindingSubjects parses a string of role binding subjects into a slice of subjects
//
// Example:
//
//	rules, err := ParseRoleBindingSubjects("kind=ServiceAccount;name=foo;namespace=bar,kind=ServiceAccount;name=foo2;namespace=bar2")
//	if err != nil {
//		// handle error
//	}
//	fmt.Println(rules) // [{Kind:ServiceAccount Name:foo Namespace:bar} {Kind:ServiceAccount Name:foo2 Namespace:bar2}]
func ParseRoleBindingSubjects(rulesStr string) ([]rbacv1.Subject, error) {
	subjects := []rbacv1.Subject{}
	rules := strc.Array(rulesStr)
	for roleIndex, rule := range rules {
		subject := rbacv1.Subject{}
		properties := strc.Properties(rule)
		for keyIndex, item := range properties {
			key, value, err := strc.KeyValue(item)
			if err != nil {
				return nil, err
			}
			switch key {
			case KeyKind:
				subject.Kind = value
			case KeyAPIGroup:
				subject.APIGroup = value
			case KeyName:
				subject.Name = value
			case KeyNamespace:
				subject.Namespace = value
			default:
				return nil, fmt.Errorf("%w at index %d in key index %d with name %q", ErrInvalidKeyInRole, roleIndex, keyIndex, key)
			}
		}
		subjects = append(subjects, subject)
	}
	return subjects, nil
}

// ParseRoleBindingRoleRef parses a string of role binding role ref into a role ref
//
// Example:
//
//	roleRef, err := ParseRoleBindingRoleRef("kind:Role;apiGroup:rbac.authorization.k8s.io;name:my-role")
//	if err != nil {
//		// handle error
//	}
//	fmt.Println(roleRef) // {APIGroup:rbac.authorization.k8s.io Kind:Role Name:my-role}
func ParseRoleBindingRoleRef(str string) (rbacv1.RoleRef, error) {
	rf := rbacv1.RoleRef{}
	properties := strc.Properties(str)
	for _, p := range properties {
		key, value, err := strc.KeyValue(p)
		if err != nil {
			return rbacv1.RoleRef{}, err
		}
		switch key {
		case KeyAPIGroup:
			rf.APIGroup = value
		case KeyKind:
			rf.Kind = value
		case KeyName:
			rf.Name = value
		default:
			return rbacv1.RoleRef{}, fmt.Errorf("%w: %q", ErrInvalidKeyInRoleRef, key)
		}
	}
	return rf, nil
}

// "verbs=get,list;apiGroups=apps,extensions;resources=deployments,replicasets"
// "verbs=get,watch;apiGroups=;resources=pods"

// ParseCustomRole parses a string of custom role rules into a slice of policy rules
//
// Example:
//
//	rules, err := ParseCustomRole("verbs=get,list;apiGroups=apps,extensions;resources=deployments,replicasets::verbs=get,watch;apiGroups=;resources=pods")
//	if err != nil {
//		// handle error
//	}
//	fmt.Println(rules) // [{Verbs:[get list] APIGroups:[apps extensions] Resources:[deployments replicasets]} {Verbs:[get watch] APIGroups:[] Resources:[pods]}]
func ParseCustomRole(str string) ([]rbacv1.PolicyRule, error) {
	rules := []rbacv1.PolicyRule{}
	strRules := strc.ArrayC(str)
	for strIdx, strRule := range strRules {
		properties := strc.Properties(strRule)
		rule := rbacv1.PolicyRule{}
		for propIdx, prop := range properties {
			key, value, err := strc.KeyValue(prop)
			if err != nil {
				return nil, err
			}
			switch key {
			case KeyVerbs:
				rule.Verbs = strc.Array(value)
			case KeyAPIGroups:
				rule.APIGroups = strc.Array(value)
			case KeyResources:
				rule.Resources = strc.Array(value)
			case KeyResourceNames:
				rule.ResourceNames = strc.Array(value)
			default:
				return nil, fmt.Errorf("%w at index %d with property index %d and key %q", ErrInvalidKeyInCustomRole, strIdx, propIdx, key)
			}
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
