package controller

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func TestLabelChecker_Create(t *testing.T) {
	type fields struct {
		ExpectedLabel string
	}
	type args struct {
		e event.CreateEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "should return true if the label is present",
			fields: fields{
				ExpectedLabel: LabelNamespacePermissionControl,
			},
			args: args{e: event.CreateEvent{
				Object: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
						Labels: map[string]string{
							LabelNamespacePermissionControl: "true",
						},
					},
				},
			}},
			want: true,
		},
		{
			name: "should return false if the label is not present",
			fields: fields{
				ExpectedLabel: LabelNamespacePermissionControl,
			},
			args: args{e: event.CreateEvent{
				Object: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
						Labels: map[string]string{
							"foo": "bar",
						},
					},
				},
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LabelChecker{
				ExpectedLabel: tt.fields.ExpectedLabel,
			}
			if got := l.Create(tt.args.e); got != tt.want {
				t.Errorf("LabelChecker.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLabelChecker_Delete(t *testing.T) {
	type fields struct {
		ExpectedLabel string
	}
	type args struct {
		e event.DeleteEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "should return true if the label is present",
			fields: fields{
				ExpectedLabel: LabelNamespacePermissionControl,
			},
			args: args{e: event.DeleteEvent{
				Object: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
						Labels: map[string]string{
							LabelNamespacePermissionControl: "true",
						},
					},
				},
			}},
			want: true,
		},
		{
			name: "should return false if the label is not present",
			fields: fields{
				ExpectedLabel: LabelNamespacePermissionControl,
			},
			args: args{e: event.DeleteEvent{
				Object: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
						Labels: map[string]string{
							"foo": "bar",
						},
					},
				},
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LabelChecker{
				ExpectedLabel: tt.fields.ExpectedLabel,
			}
			if got := l.Delete(tt.args.e); got != tt.want {
				t.Errorf("LabelChecker.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLabelChecker_Update(t *testing.T) {
	type fields struct {
		ExpectedLabel string
	}
	type args struct {
		e event.UpdateEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "should return true if the label is present",
			fields: fields{
				ExpectedLabel: LabelNamespacePermissionControl,
			},
			args: args{e: event.UpdateEvent{
				ObjectNew: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
						Labels: map[string]string{
							LabelNamespacePermissionControl: "true",
						},
					},
				},
			}},
			want: true,
		},
		{
			name: "should return false if the label is not present",
			fields: fields{
				ExpectedLabel: LabelNamespacePermissionControl,
			},
			args: args{e: event.UpdateEvent{
				ObjectNew: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
						Labels: map[string]string{
							"foo": "bar",
						},
					},
				},
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LabelChecker{
				ExpectedLabel: tt.fields.ExpectedLabel,
			}
			if got := l.Update(tt.args.e); got != tt.want {
				t.Errorf("LabelChecker.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLabelChecker_Generic(t *testing.T) {
	type fields struct {
		ExpectedLabel string
	}
	type args struct {
		e event.GenericEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "should return true if the label is present",
			fields: fields{
				ExpectedLabel: LabelNamespacePermissionControl,
			},
			args: args{e: event.GenericEvent{
				Object: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
						Labels: map[string]string{
							LabelNamespacePermissionControl: "true",
						},
					},
				},
			}},
			want: true,
		},
		{
			name: "should return false if the label is not present",
			fields: fields{
				ExpectedLabel: LabelNamespacePermissionControl,
			},
			args: args{e: event.GenericEvent{
				Object: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test",
						Labels: map[string]string{
							"foo": "bar",
						},
					},
				},
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LabelChecker{
				ExpectedLabel: tt.fields.ExpectedLabel,
			}
			if got := l.Generic(tt.args.e); got != tt.want {
				t.Errorf("LabelChecker.Generic() = %v, want %v", got, tt.want)
			}
		})
	}
}
