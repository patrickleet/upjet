// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package reference

import (
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	"github.com/crossplane/upjet/v2/pkg/config"
)

func TestSetReferenceTypes(t *testing.T) {
	resources := func(refs config.References) map[string]*config.Resource {
		return map[string]*config.Resource{
			"dummy_human":   {Kind: "Human", ShortGroup: "user", Version: "v1alpha1"},
			"dummy_machine": {Kind: "Machine", ShortGroup: "user", Version: "v1alpha1"},
			"dummy_root":    {Kind: "Root", Version: "v1beta1"},
			"dummy_grant":   {Kind: "Grant", ShortGroup: "user", Version: "v1alpha1", References: refs},
		}
	}
	cases := map[string]struct {
		rootGroup string
		refs      config.References
		want      config.References
		err       error
	}{
		"SingleTargetUnchanged": {
			rootGroup: "dummy.example.org",
			refs:      config.References{"user_id": {TerraformName: "dummy_human"}},
			want:      config.References{"user_id": {TerraformName: "dummy_human", Type: "github.com/example/provider/apis/user/v1alpha1.Human"}},
		},
		"MultiKind": {
			rootGroup: "dummy.example.org",
			refs: config.References{"user_id": {
				TerraformName: "dummy_human",
				AdditionalTargets: &[]config.ReferenceTarget{
					{TerraformName: "dummy_machine", Extractor: "E()"},
					{TerraformName: "dummy_root"},
					{Type: "github.com/other/apis/user/v1.Human", APIVersion: "user.other.example.org/v1"},
				},
			}},
			want: config.References{"user_id": {
				TerraformName: "dummy_human",
				Type:          "github.com/example/provider/apis/user/v1alpha1.Human",
				APIVersion:    "user.dummy.example.org/v1alpha1",
				AdditionalTargets: &[]config.ReferenceTarget{
					{TerraformName: "dummy_machine", Type: "github.com/example/provider/apis/user/v1alpha1.Machine", APIVersion: "user.dummy.example.org/v1alpha1", Extractor: "E()"},
					{TerraformName: "dummy_root", Type: "github.com/example/provider/apis/dummy/v1beta1.Root", APIVersion: "dummy.example.org/v1beta1"},
					{Type: "github.com/other/apis/user/v1.Human", APIVersion: "user.other.example.org/v1"},
				},
			}},
		},
		"MultiKindWithoutRootGroup": {
			refs: config.References{"user_id": {
				TerraformName:     "dummy_human",
				AdditionalTargets: &[]config.ReferenceTarget{{TerraformName: "dummy_machine"}},
			}},
			err: errors.Wrap(errors.New("cannot determine the API version of Terraform resource dummy_human: the root group is not set"), "cannot set the reference targets of dummy_grant.user_id"),
		},
		"MultiKindTypeWithoutAPIVersion": {
			rootGroup: "dummy.example.org",
			refs: config.References{"user_id": {
				TerraformName:     "dummy_human",
				AdditionalTargets: &[]config.ReferenceTarget{{Type: "Machine"}},
			}},
			err: errors.Wrap(errors.New(`reference target "Machine" (Terraform name "") must have a Type and an APIVersion`), "cannot set the reference targets of dummy_grant.user_id"),
		},
		"MultiKindUnknownTerraformName": {
			rootGroup: "dummy.example.org",
			refs: config.References{"user_id": {
				TerraformName:     "dummy_human",
				AdditionalTargets: &[]config.ReferenceTarget{{TerraformName: "dummy_unknown"}},
			}},
			err: errors.Wrap(errors.New("cannot find configuration for Terraform resource: dummy_unknown"), "cannot set the reference targets of dummy_grant.user_id"),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			rs := resources(tc.refs)
			rr := &Injector{ModulePath: "github.com/example/provider/apis", ProviderShortName: "dummy", RootGroup: tc.rootGroup}
			err := rr.SetReferenceTypes(rs)
			if diff := cmp.Diff(tc.err, err, test.EquateErrors()); diff != "" {
				t.Fatalf("SetReferenceTypes(): -want error, +got error:\n%s", diff)
			}
			if tc.err != nil {
				return
			}
			if diff := cmp.Diff(tc.want, rs["dummy_grant"].References); diff != "" {
				t.Errorf("SetReferenceTypes(): -want references, +got references:\n%s", diff)
			}
		})
	}
}
