// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package kindref

import (
	"testing"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	"github.com/google/go-cmp/cmp"
	"k8s.io/utils/ptr"
)

var policy = &xpv2.Policy{Resolve: ptr.To(xpv2.ResolvePolicyAlways)}

func TestReference(t *testing.T) {
	var none *Reference
	if none.GetAPIVersion() != "" || none.GetKind() != "" || none.ToReference() != nil {
		t.Fatal("nil Reference: want empty apiVersion and kind, and a nil reference")
	}
	r := &Reference{APIVersion: "g.example.org/v1", Kind: "Machine", Name: "a", Policy: policy}
	if r.GetAPIVersion() != "g.example.org/v1" || r.GetKind() != "Machine" {
		t.Errorf("GetAPIVersion(), GetKind(): got %q, %q", r.GetAPIVersion(), r.GetKind())
	}
	if diff := cmp.Diff(&xpv2.Reference{Name: "a", Policy: policy}, r.ToReference()); diff != "" {
		t.Errorf("ToReference(): -want, +got:\n%s", diff)
	}
	if diff := cmp.Diff(r, NewReference("g.example.org/v1", "Machine", r.ToReference())); diff != "" {
		t.Errorf("NewReference(ToReference()): -want, +got:\n%s", diff)
	}
	if NewReference("g.example.org/v1", "Machine", nil) != nil {
		t.Error("NewReference(nil): want nil")
	}
}

func TestNamespacedReference(t *testing.T) {
	var none *NamespacedReference
	if none.GetAPIVersion() != "" || none.GetKind() != "" || none.ToReference() != nil {
		t.Fatal("nil NamespacedReference: want empty apiVersion and kind, and a nil reference")
	}
	r := &NamespacedReference{APIVersion: "g.example.org/v1", Kind: "Machine", Name: "a", Namespace: "ns", Policy: policy}
	if r.GetAPIVersion() != "g.example.org/v1" || r.GetKind() != "Machine" {
		t.Errorf("GetAPIVersion(), GetKind(): got %q, %q", r.GetAPIVersion(), r.GetKind())
	}
	if diff := cmp.Diff(&xpv2.NamespacedReference{Name: "a", Namespace: "ns", Policy: policy}, r.ToReference()); diff != "" {
		t.Errorf("ToReference(): -want, +got:\n%s", diff)
	}
	if diff := cmp.Diff(r, NewNamespacedReference("g.example.org/v1", "Machine", r.ToReference())); diff != "" {
		t.Errorf("NewNamespacedReference(ToReference()): -want, +got:\n%s", diff)
	}
	if NewNamespacedReference("g.example.org/v1", "Machine", nil) != nil {
		t.Error("NewNamespacedReference(nil): want nil")
	}
}

func TestSelector(t *testing.T) {
	var none *Selector
	if none.GetAPIVersion() != "" || none.GetKind() != "" || none.ToSelector() != nil {
		t.Fatal("nil Selector: want empty apiVersion and kind, and a nil selector")
	}
	s := &Selector{APIVersion: "g.example.org/v1", Kind: "Machine", MatchLabels: map[string]string{"a": "b"}, MatchControllerRef: ptr.To(true), Policy: policy}
	if s.GetAPIVersion() != "g.example.org/v1" || s.GetKind() != "Machine" {
		t.Errorf("GetAPIVersion(), GetKind(): got %q, %q", s.GetAPIVersion(), s.GetKind())
	}
	want := &xpv2.Selector{MatchLabels: map[string]string{"a": "b"}, MatchControllerRef: ptr.To(true), Policy: policy}
	if diff := cmp.Diff(want, s.ToSelector()); diff != "" {
		t.Errorf("ToSelector(): -want, +got:\n%s", diff)
	}
}

func TestNamespacedSelector(t *testing.T) {
	var none *NamespacedSelector
	if none.GetAPIVersion() != "" || none.GetKind() != "" || none.ToSelector() != nil {
		t.Fatal("nil NamespacedSelector: want empty apiVersion and kind, and a nil selector")
	}
	s := &NamespacedSelector{APIVersion: "g.example.org/v1", Kind: "Machine", MatchLabels: map[string]string{"a": "b"}, MatchControllerRef: ptr.To(true), Policy: policy, Namespace: "ns"}
	if s.GetAPIVersion() != "g.example.org/v1" || s.GetKind() != "Machine" {
		t.Errorf("GetAPIVersion(), GetKind(): got %q, %q", s.GetAPIVersion(), s.GetKind())
	}
	want := &xpv2.NamespacedSelector{MatchLabels: map[string]string{"a": "b"}, MatchControllerRef: ptr.To(true), Policy: policy, Namespace: "ns"}
	if diff := cmp.Diff(want, s.ToSelector()); diff != "" {
		t.Errorf("ToSelector(): -want, +got:\n%s", diff)
	}
}
