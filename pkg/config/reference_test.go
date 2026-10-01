// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

func TestReferenceTargets(t *testing.T) {
	if got := (Reference{TerraformName: "a"}).Targets(); got != nil {
		t.Errorf("Targets() of a single-target reference: want nil, got %v", got)
	}
	r := Reference{
		Type:              "github.com/example/apis/user/v1.Human",
		TerraformName:     "human",
		APIVersion:        "user.example.org/v1",
		Extractor:         "E()",
		RefFieldName:      "UserRef",
		AdditionalTargets: []ReferenceTarget{{Type: "Machine", APIVersion: "user.example.org/v1"}},
	}
	want := []ReferenceTarget{
		{Type: "github.com/example/apis/user/v1.Human", TerraformName: "human", APIVersion: "user.example.org/v1", Extractor: "E()"},
		{Type: "Machine", APIVersion: "user.example.org/v1"},
	}
	if diff := cmp.Diff(want, r.Targets()); diff != "" {
		t.Errorf("Targets(): -want, +got:\n%s", diff)
	}
	if got := want[0].Kind(); got != "Human" {
		t.Errorf("Kind() of a qualified type: want Human, got %q", got)
	}
	if got := want[1].Kind(); got != "Machine" {
		t.Errorf("Kind() of a local type: want Machine, got %q", got)
	}
}

func TestIsAmbiguousKind(t *testing.T) {
	targets := []ReferenceTarget{
		{Type: "a/project/v1.Grant", APIVersion: "project.example.org/v1"},
		{Type: "a/user/v1.Grant", APIVersion: "user.example.org/v1"},
		{Type: "a/user/v1.Machine", APIVersion: "user.example.org/v1"},
	}
	if !IsAmbiguousKind(targets, "Grant") {
		t.Error("IsAmbiguousKind(Grant): want true")
	}
	if IsAmbiguousKind(targets, "Machine") {
		t.Error("IsAmbiguousKind(Machine): want false")
	}
	if IsAmbiguousKind(targets, "Other") {
		t.Error("IsAmbiguousKind(Other): want false")
	}
}

func TestReferenceValidateTargets(t *testing.T) {
	cases := map[string]struct {
		ref  Reference
		want error
	}{
		"SingleTarget": {
			ref: Reference{TerraformName: "a"},
		},
		"Valid": {
			ref: Reference{
				Type: "Human", APIVersion: "user.example.org/v1",
				AdditionalTargets: []ReferenceTarget{{Type: "Machine", APIVersion: "user.example.org/v1"}},
			},
		},
		"SameKindDifferentGroups": {
			ref: Reference{
				Type: "a/project/v1.Grant", APIVersion: "project.example.org/v1",
				AdditionalTargets: []ReferenceTarget{{Type: "a/user/v1.Grant", APIVersion: "user.example.org/v1"}},
			},
		},
		"MissingAPIVersion": {
			ref: Reference{
				Type:              "Human",
				AdditionalTargets: []ReferenceTarget{{Type: "Machine", APIVersion: "user.example.org/v1"}},
			},
			want: errors.New(`reference target "Human" (Terraform name "") must have a Type and an APIVersion`),
		},
		"MissingType": {
			ref: Reference{
				Type: "Human", APIVersion: "user.example.org/v1",
				AdditionalTargets: []ReferenceTarget{{TerraformName: "machine"}},
			},
			want: errors.New(`reference target "" (Terraform name "machine") must have a Type and an APIVersion`),
		},
		"MalformedAPIVersion": {
			ref: Reference{
				Type: "Human", APIVersion: "v1",
				AdditionalTargets: []ReferenceTarget{{Type: "Machine", APIVersion: "user.example.org/v1"}},
			},
			want: errors.New(`reference target "Human" must have an APIVersion of the form group/version, got "v1"`),
		},
		"Duplicate": {
			ref: Reference{
				Type: "a/user/v1.Human", APIVersion: "user.example.org/v1",
				AdditionalTargets: []ReferenceTarget{{Type: "b/user/v1.Human", APIVersion: "user.example.org/v1"}},
			},
			want: errors.New(`reference targets must be distinct, but user.example.org/v1, Kind=Human is configured more than once`),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if diff := cmp.Diff(tc.want, tc.ref.ValidateTargets(), test.EquateErrors()); diff != "" {
				t.Errorf("ValidateTargets(): -want error, +got error:\n%s", diff)
			}
		})
	}
}
