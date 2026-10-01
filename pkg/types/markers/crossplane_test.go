// SPDX-FileCopyrightText: 2023 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package markers

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/crossplane/upjet/v2/pkg/config"
)

func TestCrossplaneOptions_String(t *testing.T) {
	type args struct {
		referenceToType            string
		referenceExtractor         string
		referenceFieldName         string
		referenceSelectorFieldName string
	}
	type want struct {
		out string
	}
	cases := map[string]struct {
		args
		want
	}{
		"NoOption": {
			args: args{
				referenceToType: "",
			},
			want: want{
				out: "",
			},
		},
		"WithType": {
			args: args{
				referenceToType: "SecurityGroup",
			},
			want: want{
				out: "+crossplane:generate:reference:type=SecurityGroup\n",
			},
		},
		"WithAll": {
			args: args{
				referenceToType:            "github.com/crossplane/provider-aws/apis/ec2/v1beta1.Subnet",
				referenceExtractor:         "github.com/crossplane/provider-aws/apis/ec2/v1beta1.SubnetARN()",
				referenceFieldName:         "SubnetIDRefs",
				referenceSelectorFieldName: "SubnetIDSelector",
			},
			want: want{
				out: `+crossplane:generate:reference:type=github.com/crossplane/provider-aws/apis/ec2/v1beta1.Subnet
+crossplane:generate:reference:extractor=github.com/crossplane/provider-aws/apis/ec2/v1beta1.SubnetARN()
+crossplane:generate:reference:refFieldName=SubnetIDRefs
+crossplane:generate:reference:selectorFieldName=SubnetIDSelector
`,
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			o := CrossplaneOptions{
				Reference: config.Reference{
					Type:              tc.referenceToType,
					Extractor:         tc.referenceExtractor,
					RefFieldName:      tc.referenceFieldName,
					SelectorFieldName: tc.referenceSelectorFieldName,
				},
			}
			got := o.String()
			if diff := cmp.Diff(tc.want.out, got); diff != "" {
				t.Errorf("CrossplaneOptions.String(): -want result, +got result: %s", diff)
			}
		})
	}
}

func TestCrossplaneOptions_StringMultiKind(t *testing.T) {
	cases := map[string]struct {
		ref  config.Reference
		want string
	}{
		"DefaultExtractors": {
			ref: config.Reference{
				Type:              "github.com/example/provider/apis/user/v1alpha1.Human",
				APIVersion:        "user.example.org/v1alpha1",
				RefFieldName:      "UserRef",
				AdditionalTargets: []config.ReferenceTarget{{Type: "Machine", APIVersion: "user.example.org/v1alpha1"}},
			},
			want: `+crossplane:generate:reference:type=github.com/example/provider/apis/user/v1alpha1.Human
+crossplane:generate:reference:apiVersion=user.example.org/v1alpha1
+crossplane:generate:reference:type=Machine
+crossplane:generate:reference:apiVersion=user.example.org/v1alpha1
+crossplane:generate:reference:refFieldName=UserRef
`,
		},
		"OneCustomExtractor": {
			ref: config.Reference{
				Type:       "Human",
				APIVersion: "user.example.org/v1alpha1",
				AdditionalTargets: []config.ReferenceTarget{{
					Type:       "Machine",
					APIVersion: "user.example.org/v1alpha1",
					Extractor:  "github.com/example/provider/config/common.MachineID()",
				}},
			},
			want: `+crossplane:generate:reference:type=Human
+crossplane:generate:reference:apiVersion=user.example.org/v1alpha1
+crossplane:generate:reference:extractor=github.com/crossplane/crossplane-runtime/v2/pkg/reference.ExternalName()
+crossplane:generate:reference:type=Machine
+crossplane:generate:reference:apiVersion=user.example.org/v1alpha1
+crossplane:generate:reference:extractor=github.com/example/provider/config/common.MachineID()
`,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := CrossplaneOptions{Reference: tc.ref}.String()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("CrossplaneOptions.String(): -want result, +got result: %s", diff)
			}
		})
	}
}
