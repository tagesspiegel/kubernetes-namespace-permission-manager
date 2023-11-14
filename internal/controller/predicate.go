package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var (
	_ predicate.Predicate = &LabelChecker{}
)

type LabelChecker struct {
	ExpectedLabel string
}

func (l *LabelChecker) Create(e event.CreateEvent) bool {
	_, ok := e.Object.GetLabels()[l.ExpectedLabel]
	return ok
}

func (l *LabelChecker) Delete(e event.DeleteEvent) bool {
	_, ok := e.Object.GetLabels()[l.ExpectedLabel]
	return ok
}

func (l *LabelChecker) Update(e event.UpdateEvent) bool {
	_, ok := e.ObjectNew.GetLabels()[l.ExpectedLabel]
	return ok
}

func (l *LabelChecker) Generic(e event.GenericEvent) bool {
	_, ok := e.Object.GetLabels()[l.ExpectedLabel]
	return ok
}
