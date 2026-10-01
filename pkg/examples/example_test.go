// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package examples

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/crossplane/upjet/v2/pkg/config"
)

func TestGetSelectorField(t *testing.T) {
	multiKind := config.Reference{
		TerraformName: "dummy_human",
		Type:          "a/user/v1.Human",
		APIVersion:    "user.example.org/v1",
		AdditionalTargets: &[]config.ReferenceTarget{
			{TerraformName: "dummy_machine", Type: "a/user/v1.Machine", APIVersion: "user.example.org/v1"},
			{TerraformName: "dummy_project_grant", Type: "a/project/v1.Grant", APIVersion: "project.example.org/v1"},
			{TerraformName: "dummy_user_grant", Type: "a/user/v1.Grant", APIVersion: "user.example.org/v1"},
		},
	}
	labels := func(name string) map[string]string {
		return map[string]string{labelExampleName: name}
	}
	cases := map[string]struct {
		value any
		cfg   config.Reference
		want  any
	}{
		"SingleTarget": {
			value: "${dummy_human.alice.id}",
			cfg:   config.Reference{TerraformName: "dummy_human"},
			want:  map[string]any{"matchLabels": labels("alice")},
		},
		"MultiKindDefaultTarget": {
			value: "${dummy_human.alice.id}",
			cfg:   multiKind,
			want:  map[string]any{"matchLabels": labels("alice")},
		},
		"MultiKindOtherTarget": {
			value: "${dummy_machine.bot.id}",
			cfg:   multiKind,
			want:  map[string]any{"matchLabels": labels("bot"), "kind": "Machine"},
		},
		"MultiKindAmbiguousKind": {
			value: "${dummy_user_grant.g.id}",
			cfg:   multiKind,
			want:  map[string]any{"matchLabels": labels("g"), "kind": "Grant", "apiVersion": "user.example.org/v1"},
		},
		// The default and an additional target can be the same Terraform
		// resource at different API versions; the default needs no kind.
		"SameTerraformResourceAsDefault": {
			value: "${dummy_human.alice.id}",
			cfg: config.Reference{
				TerraformName: "dummy_human",
				Type:          "a/user/v1.Human",
				APIVersion:    "user.example.org/v1",
				AdditionalTargets: &[]config.ReferenceTarget{
					{TerraformName: "dummy_human", Type: "a/user/v2.Human", APIVersion: "user.example.org/v2"},
				},
			},
			want: map[string]any{"matchLabels": labels("alice")},
		},
		"NotAReference": {
			value: "literal",
			cfg:   multiKind,
			want:  map[string]any{"matchLabels": labels(defaultExampleName)},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if diff := cmp.Diff(tc.want, getSelectorField(tc.value, tc.cfg)); diff != "" {
				t.Errorf("getSelectorField(): -want, +got:\n%s", diff)
			}
		})
	}
}
