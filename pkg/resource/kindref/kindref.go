// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

// Package kindref contains the reference and selector types generated for
// fields that can be resolved from more than one kind of managed resource.
//
// A field configured with config.Reference.AdditionalTargets gets a
// <field>Ref of type Reference (NamespacedReference for namespaced managed
// resources) and a <field>Selector of type Selector (NamespacedSelector).
// These are the standard Crossplane reference and selector, plus an optional
// apiVersion and kind that choose which of the field's configured targets to
// resolve. Omitting both resolves the field's default (first) target.
//
// The generated ResolveReferences only ever looks up the configured
// targets, so these types never trigger discovery of arbitrary kinds.
// +kubebuilder:object:generate=true
package kindref

import (
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

// A Reference to a named object of one of the kinds a field accepts.
type Reference struct {
	// APIVersion of the referenced object, as group/version. Only needed
	// when Kind matches more than one of the field's targets.
	// +optional
	// +kubebuilder:validation:MinLength=1
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the referenced object. Defaults to the field's default kind.
	// +optional
	// +kubebuilder:validation:MinLength=1
	Kind string `json:"kind,omitempty"`

	// Name of the referenced object.
	Name string `json:"name"`

	// Policies for referencing.
	// +optional
	Policy *xpv2.Policy `json:"policy,omitempty"`
}

// A NamespacedReference to a named object of one of the kinds a field
// accepts.
type NamespacedReference struct {
	// APIVersion of the referenced object, as group/version. Only needed
	// when Kind matches more than one of the field's targets.
	// +optional
	// +kubebuilder:validation:MinLength=1
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the referenced object. Defaults to the field's default kind.
	// +optional
	// +kubebuilder:validation:MinLength=1
	Kind string `json:"kind,omitempty"`

	// Name of the referenced object.
	Name string `json:"name"`

	// Namespace of the referenced object. Defaults to the namespace of the
	// referencing object.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Policies for referencing.
	// +optional
	Policy *xpv2.Policy `json:"policy,omitempty"`
}

// A Selector selects an object of one of the kinds a field accepts.
type Selector struct {
	// APIVersion of the selected object, as group/version. Only needed when
	// Kind matches more than one of the field's targets.
	// +optional
	// +kubebuilder:validation:MinLength=1
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the selected object. Defaults to the field's default kind.
	// +optional
	// +kubebuilder:validation:MinLength=1
	Kind string `json:"kind,omitempty"`

	// MatchLabels ensures an object with matching labels is selected.
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// MatchControllerRef ensures an object with the same controller reference
	// as the selecting object is selected.
	MatchControllerRef *bool `json:"matchControllerRef,omitempty"`

	// Policies for selection.
	// +optional
	Policy *xpv2.Policy `json:"policy,omitempty"`
}

// A NamespacedSelector selects a namespaced object of one of the kinds a
// field accepts.
type NamespacedSelector struct {
	// APIVersion of the selected object, as group/version. Only needed when
	// Kind matches more than one of the field's targets.
	// +optional
	// +kubebuilder:validation:MinLength=1
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the selected object. Defaults to the field's default kind.
	// +optional
	// +kubebuilder:validation:MinLength=1
	Kind string `json:"kind,omitempty"`

	// MatchLabels ensures an object with matching labels is selected.
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// MatchControllerRef ensures an object with the same controller reference
	// as the selecting object is selected.
	MatchControllerRef *bool `json:"matchControllerRef,omitempty"`

	// Policies for selection.
	// +optional
	Policy *xpv2.Policy `json:"policy,omitempty"`

	// Namespace for the selector. Defaults to the namespace of the
	// referencing object.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// GetAPIVersion returns the referenced apiVersion, or "" if r is nil.
func (r *Reference) GetAPIVersion() string {
	if r == nil {
		return ""
	}
	return r.APIVersion
}

// GetKind returns the referenced kind, or "" if r is nil.
func (r *Reference) GetKind() string {
	if r == nil {
		return ""
	}
	return r.Kind
}

// ToReference returns the standard reference to resolve, or nil if r is nil.
func (r *Reference) ToReference() *xpv2.Reference {
	if r == nil {
		return nil
	}
	return &xpv2.Reference{Name: r.Name, Policy: r.Policy}
}

// NewReference returns the Reference to store for a resolved reference to
// the supplied apiVersion and kind. It returns nil if resolved is nil.
func NewReference(apiVersion, kind string, resolved *xpv2.Reference) *Reference {
	if resolved == nil {
		return nil
	}
	return &Reference{APIVersion: apiVersion, Kind: kind, Name: resolved.Name, Policy: resolved.Policy}
}

// GetAPIVersion returns the referenced apiVersion, or "" if r is nil.
func (r *NamespacedReference) GetAPIVersion() string {
	if r == nil {
		return ""
	}
	return r.APIVersion
}

// GetKind returns the referenced kind, or "" if r is nil.
func (r *NamespacedReference) GetKind() string {
	if r == nil {
		return ""
	}
	return r.Kind
}

// ToReference returns the standard reference to resolve, or nil if r is nil.
func (r *NamespacedReference) ToReference() *xpv2.NamespacedReference {
	if r == nil {
		return nil
	}
	return &xpv2.NamespacedReference{Name: r.Name, Namespace: r.Namespace, Policy: r.Policy}
}

// NewNamespacedReference returns the NamespacedReference to store for a
// resolved reference to the supplied apiVersion and kind. It returns nil if
// resolved is nil.
func NewNamespacedReference(apiVersion, kind string, resolved *xpv2.NamespacedReference) *NamespacedReference {
	if resolved == nil {
		return nil
	}
	return &NamespacedReference{APIVersion: apiVersion, Kind: kind, Name: resolved.Name, Namespace: resolved.Namespace, Policy: resolved.Policy}
}

// GetAPIVersion returns the selected apiVersion, or "" if s is nil.
func (s *Selector) GetAPIVersion() string {
	if s == nil {
		return ""
	}
	return s.APIVersion
}

// GetKind returns the selected kind, or "" if s is nil.
func (s *Selector) GetKind() string {
	if s == nil {
		return ""
	}
	return s.Kind
}

// ToSelector returns the standard selector to resolve, or nil if s is nil.
func (s *Selector) ToSelector() *xpv2.Selector {
	if s == nil {
		return nil
	}
	return &xpv2.Selector{MatchLabels: s.MatchLabels, MatchControllerRef: s.MatchControllerRef, Policy: s.Policy}
}

// GetAPIVersion returns the selected apiVersion, or "" if s is nil.
func (s *NamespacedSelector) GetAPIVersion() string {
	if s == nil {
		return ""
	}
	return s.APIVersion
}

// GetKind returns the selected kind, or "" if s is nil.
func (s *NamespacedSelector) GetKind() string {
	if s == nil {
		return ""
	}
	return s.Kind
}

// ToSelector returns the standard selector to resolve, or nil if s is nil.
func (s *NamespacedSelector) ToSelector() *xpv2.NamespacedSelector {
	if s == nil {
		return nil
	}
	return &xpv2.NamespacedSelector{MatchLabels: s.MatchLabels, MatchControllerRef: s.MatchControllerRef, Policy: s.Policy, Namespace: s.Namespace}
}
