// SPDX-FileCopyrightText: 2023 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/google/go-cmp/cmp"
	twtypes "github.com/muvaf/typewriter/pkg/types"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/types/name"
)

func TestBuilder_generateReferenceFields(t *testing.T) {
	tp := types.NewPackage("github.com/crossplane/upjet/v2/pkg/types", "tjtypes")

	type args struct {
		t        *types.TypeName
		f        *Field
		crdScope CRDScope
	}
	type want struct {
		outFields   []*types.Var
		outTags     []string
		outComments twtypes.Comments
	}
	cases := map[string]struct {
		args
		want
	}{
		"OnlyRefType": {
			args: args{
				crdScope: CRDScopeCluster,
				t:        types.NewTypeName(token.NoPos, tp, "Params", types.Universe.Lookup("string").Type()),
				f: &Field{
					Name: name.NewFromCamel("TestField"),
					Reference: &config.Reference{
						Type: "testObject",
					},
					FieldType: types.Universe.Lookup("string").Type(),
				},
			}, want: want{
				outFields: []*types.Var{
					types.NewField(token.NoPos, tp, "TestFieldRef", types.NewPointer(typeReferenceField), false),
					types.NewField(token.NoPos, tp, "TestFieldSelector", types.NewPointer(typeSelectorField), false),
				},
				outTags: []string{
					`json:"testFieldRef,omitempty" tf:"-"`,
					`json:"testFieldSelector,omitempty" tf:"-"`,
				},
				outComments: twtypes.Comments{
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldRef":      "// Reference to a testObject to populate testField.\n// +kubebuilder:validation:Optional\n",
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldSelector": "// Selector for a testObject to populate testField.\n// +kubebuilder:validation:Optional\n",
				},
			},
		},
		"OnlyRefTypeSlice": {
			args: args{
				crdScope: CRDScopeCluster,
				t:        types.NewTypeName(token.NoPos, tp, "Params", types.Universe.Lookup("string").Type()),
				f: &Field{
					Name: name.NewFromCamel("TestField"),
					Reference: &config.Reference{
						Type: "testObject",
					},
					FieldType: types.NewSlice(types.Universe.Lookup("string").Type()),
				},
			}, want: want{
				outFields: []*types.Var{
					types.NewField(token.NoPos, tp, "TestFieldRefs", types.NewSlice(typeReferenceField), false),
					types.NewField(token.NoPos, tp, "TestFieldSelector", types.NewPointer(typeSelectorField), false),
				},
				outTags: []string{
					`json:"testFieldRefs,omitempty" tf:"-"`,
					`json:"testFieldSelector,omitempty" tf:"-"`,
				},
				outComments: twtypes.Comments{
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldRefs":     "// References to testObject to populate testField.\n// +kubebuilder:validation:Optional\n",
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldSelector": "// Selector for a list of testObject to populate testField.\n// +kubebuilder:validation:Optional\n",
				},
			},
		},
		"WithCustomFieldName": {
			args: args{
				crdScope: CRDScopeCluster,
				t:        types.NewTypeName(token.NoPos, tp, "Params", types.Universe.Lookup("string").Type()),
				f: &Field{
					Name: name.NewFromCamel("TestField"),
					Reference: &config.Reference{
						Type:         "TestObject",
						RefFieldName: "CustomRef",
					},
					FieldType: types.Universe.Lookup("string").Type(),
				},
			}, want: want{
				outFields: []*types.Var{
					types.NewField(token.NoPos, tp, "CustomRef", types.NewPointer(typeReferenceField), false),
					types.NewField(token.NoPos, tp, "TestFieldSelector", types.NewPointer(typeSelectorField), false),
				},
				outTags: []string{
					`json:"customRef,omitempty" tf:"-"`,
					`json:"testFieldSelector,omitempty" tf:"-"`,
				},
				outComments: twtypes.Comments{
					"github.com/crossplane/upjet/v2/pkg/types.Params:CustomRef":         "// Reference to a TestObject to populate testField.\n// +kubebuilder:validation:Optional\n",
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldSelector": "// Selector for a TestObject to populate testField.\n// +kubebuilder:validation:Optional\n",
				},
			},
		},
		"WithCustomSelectorName": {
			args: args{
				crdScope: CRDScopeCluster,
				t:        types.NewTypeName(token.NoPos, tp, "Params", types.Universe.Lookup("string").Type()),
				f: &Field{
					Name: name.NewFromCamel("TestField"),
					Reference: &config.Reference{
						Type:              "TestObject",
						SelectorFieldName: "CustomSelector",
					},
					FieldType: types.Universe.Lookup("string").Type(),
				},
			}, want: want{
				outFields: []*types.Var{
					types.NewField(token.NoPos, tp, "TestFieldRef", types.NewPointer(typeReferenceField), false),
					types.NewField(token.NoPos, tp, "CustomSelector", types.NewPointer(typeSelectorField), false),
				},
				outTags: []string{
					`json:"testFieldRef,omitempty" tf:"-"`,
					`json:"customSelector,omitempty" tf:"-"`,
				},
				outComments: twtypes.Comments{
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldRef":   "// Reference to a TestObject to populate testField.\n// +kubebuilder:validation:Optional\n",
					"github.com/crossplane/upjet/v2/pkg/types.Params:CustomSelector": "// Selector for a TestObject to populate testField.\n// +kubebuilder:validation:Optional\n",
				},
			},
		},
		"ReferenceToAnotherPackage": {
			args: args{
				crdScope: CRDScopeCluster,
				t:        types.NewTypeName(token.NoPos, tp, "Params", types.Universe.Lookup("string").Type()),
				f: &Field{
					Name: name.NewFromCamel("TestField"),
					Reference: &config.Reference{
						Type: "github.com/upbound/official-providers/provider-aws/apis/somepackage/v1beta1.TestObject",
					},
					FieldType: types.Universe.Lookup("string").Type(),
				},
			}, want: want{
				outFields: []*types.Var{
					types.NewField(token.NoPos, tp, "TestFieldRef", types.NewPointer(typeReferenceField), false),
					types.NewField(token.NoPos, tp, "TestFieldSelector", types.NewPointer(typeSelectorField), false),
				},
				outTags: []string{
					`json:"testFieldRef,omitempty" tf:"-"`,
					`json:"testFieldSelector,omitempty" tf:"-"`,
				},
				outComments: twtypes.Comments{
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldRef":      "// Reference to a TestObject in somepackage to populate testField.\n// +kubebuilder:validation:Optional\n",
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldSelector": "// Selector for a TestObject in somepackage to populate testField.\n// +kubebuilder:validation:Optional\n",
				},
			},
		},
		// namespaced CRD tests
		"OnlyRefType_namespaced": {
			args: args{
				crdScope: CRDScopeNamespaced,
				t:        types.NewTypeName(token.NoPos, tp, "Params", types.Universe.Lookup("string").Type()),
				f: &Field{
					Name: name.NewFromCamel("TestField"),
					Reference: &config.Reference{
						Type: "testObject",
					},
					FieldType: types.Universe.Lookup("string").Type(),
				},
			}, want: want{
				outFields: []*types.Var{
					types.NewField(token.NoPos, tp, "TestFieldRef", types.NewPointer(typeNamespacedReferenceField), false),
					types.NewField(token.NoPos, tp, "TestFieldSelector", types.NewPointer(typeNamespacedSelectorField), false),
				},
				outTags: []string{
					`json:"testFieldRef,omitempty" tf:"-"`,
					`json:"testFieldSelector,omitempty" tf:"-"`,
				},
				outComments: twtypes.Comments{
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldRef":      "// Reference to a testObject to populate testField.\n// +kubebuilder:validation:Optional\n",
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldSelector": "// Selector for a testObject to populate testField.\n// +kubebuilder:validation:Optional\n",
				},
			},
		},
		"OnlyRefTypeSlice_namespaced": {
			args: args{
				crdScope: CRDScopeNamespaced,
				t:        types.NewTypeName(token.NoPos, tp, "Params", types.Universe.Lookup("string").Type()),
				f: &Field{
					Name: name.NewFromCamel("TestField"),
					Reference: &config.Reference{
						Type: "testObject",
					},
					FieldType: types.NewSlice(types.Universe.Lookup("string").Type()),
				},
			}, want: want{
				outFields: []*types.Var{
					types.NewField(token.NoPos, tp, "TestFieldRefs", types.NewSlice(typeNamespacedReferenceField), false),
					types.NewField(token.NoPos, tp, "TestFieldSelector", types.NewPointer(typeNamespacedSelectorField), false),
				},
				outTags: []string{
					`json:"testFieldRefs,omitempty" tf:"-"`,
					`json:"testFieldSelector,omitempty" tf:"-"`,
				},
				outComments: twtypes.Comments{
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldRefs":     "// References to testObject to populate testField.\n// +kubebuilder:validation:Optional\n",
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldSelector": "// Selector for a list of testObject to populate testField.\n// +kubebuilder:validation:Optional\n",
				},
			},
		},
		"WithCustomFieldName_namespaced": {
			args: args{
				crdScope: CRDScopeNamespaced,
				t:        types.NewTypeName(token.NoPos, tp, "Params", types.Universe.Lookup("string").Type()),
				f: &Field{
					Name: name.NewFromCamel("TestField"),
					Reference: &config.Reference{
						Type:         "TestObject",
						RefFieldName: "CustomRef",
					},
					FieldType: types.Universe.Lookup("string").Type(),
				},
			}, want: want{
				outFields: []*types.Var{
					types.NewField(token.NoPos, tp, "CustomRef", types.NewPointer(typeNamespacedReferenceField), false),
					types.NewField(token.NoPos, tp, "TestFieldSelector", types.NewPointer(typeNamespacedSelectorField), false),
				},
				outTags: []string{
					`json:"customRef,omitempty" tf:"-"`,
					`json:"testFieldSelector,omitempty" tf:"-"`,
				},
				outComments: twtypes.Comments{
					"github.com/crossplane/upjet/v2/pkg/types.Params:CustomRef":         "// Reference to a TestObject to populate testField.\n// +kubebuilder:validation:Optional\n",
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldSelector": "// Selector for a TestObject to populate testField.\n// +kubebuilder:validation:Optional\n",
				},
			},
		},
		"WithCustomSelectorName_namespaced": {
			args: args{
				crdScope: CRDScopeNamespaced,
				t:        types.NewTypeName(token.NoPos, tp, "Params", types.Universe.Lookup("string").Type()),
				f: &Field{
					Name: name.NewFromCamel("TestField"),
					Reference: &config.Reference{
						Type:              "TestObject",
						SelectorFieldName: "CustomSelector",
					},
					FieldType: types.Universe.Lookup("string").Type(),
				},
			}, want: want{
				outFields: []*types.Var{
					types.NewField(token.NoPos, tp, "TestFieldRef", types.NewPointer(typeNamespacedReferenceField), false),
					types.NewField(token.NoPos, tp, "CustomSelector", types.NewPointer(typeNamespacedSelectorField), false),
				},
				outTags: []string{
					`json:"testFieldRef,omitempty" tf:"-"`,
					`json:"customSelector,omitempty" tf:"-"`,
				},
				outComments: twtypes.Comments{
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldRef":   "// Reference to a TestObject to populate testField.\n// +kubebuilder:validation:Optional\n",
					"github.com/crossplane/upjet/v2/pkg/types.Params:CustomSelector": "// Selector for a TestObject to populate testField.\n// +kubebuilder:validation:Optional\n",
				},
			},
		},
		"ReferenceToAnotherPackage_namespaced": {
			args: args{
				crdScope: CRDScopeNamespaced,
				t:        types.NewTypeName(token.NoPos, tp, "Params", types.Universe.Lookup("string").Type()),
				f: &Field{
					Name: name.NewFromCamel("TestField"),
					Reference: &config.Reference{
						Type: "github.com/upbound/official-providers/provider-aws/apis/somepackage/v1beta1.TestObject",
					},
					FieldType: types.Universe.Lookup("string").Type(),
				},
			}, want: want{
				outFields: []*types.Var{
					types.NewField(token.NoPos, tp, "TestFieldRef", types.NewPointer(typeNamespacedReferenceField), false),
					types.NewField(token.NoPos, tp, "TestFieldSelector", types.NewPointer(typeNamespacedSelectorField), false),
				},
				outTags: []string{
					`json:"testFieldRef,omitempty" tf:"-"`,
					`json:"testFieldSelector,omitempty" tf:"-"`,
				},
				outComments: twtypes.Comments{
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldRef":      "// Reference to a TestObject in somepackage to populate testField.\n// +kubebuilder:validation:Optional\n",
					"github.com/crossplane/upjet/v2/pkg/types.Params:TestFieldSelector": "// Selector for a TestObject in somepackage to populate testField.\n// +kubebuilder:validation:Optional\n",
				},
			},
		},
		"MultiKind": {
			args: args{
				crdScope: CRDScopeCluster,
				t:        types.NewTypeName(token.NoPos, tp, "Params", types.Universe.Lookup("string").Type()),
				f: &Field{
					Name: name.NewFromCamel("UserID"),
					Reference: &config.Reference{
						Type:       "github.com/upbound/official-providers/provider-dummy/apis/user/v1alpha1.Human",
						APIVersion: "user.dummy.example.org/v1alpha1",
						AdditionalTargets: []config.ReferenceTarget{{
							Type:       "github.com/upbound/official-providers/provider-dummy/apis/user/v1alpha1.Machine",
							APIVersion: "user.dummy.example.org/v1alpha1",
						}},
					},
					FieldType: types.NewPointer(types.Universe.Lookup("string").Type()),
				},
			}, want: want{
				outFields: []*types.Var{
					types.NewField(token.NoPos, tp, "UserIDRef", types.NewPointer(typeKindReferenceField), false),
					types.NewField(token.NoPos, tp, "UserIDSelector", types.NewPointer(typeKindSelectorField), false),
				},
				outTags: []string{
					`json:"userIdRef,omitempty" tf:"-"`,
					`json:"userIdSelector,omitempty" tf:"-"`,
				},
				outComments: twtypes.Comments{
					"github.com/crossplane/upjet/v2/pkg/types.Params:UserIDRef": `// Reference to a Human in user (default) or Machine in user to populate userId.
// Set kind, and apiVersion if kind alone is ambiguous, to choose a target other than the default.
// +kubebuilder:validation:XValidation:rule="!has(self.kind) || self.kind in ['Human', 'Machine']",message="kind must be one of Human, Machine"
// +kubebuilder:validation:XValidation:rule="!has(self.apiVersion) || (has(self.kind) && (self.apiVersion + '/' + self.kind) in ['user.dummy.example.org/v1alpha1/Human', 'user.dummy.example.org/v1alpha1/Machine'])",message="apiVersion and kind must be one of: user.dummy.example.org/v1alpha1 Human, user.dummy.example.org/v1alpha1 Machine"
// +kubebuilder:validation:Optional
`,
					"github.com/crossplane/upjet/v2/pkg/types.Params:UserIDSelector": `// Selector for a Human in user (default) or Machine in user to populate userId.
// Set kind, and apiVersion if kind alone is ambiguous, to choose a target other than the default.
// +kubebuilder:validation:XValidation:rule="!has(self.kind) || self.kind in ['Human', 'Machine']",message="kind must be one of Human, Machine"
// +kubebuilder:validation:XValidation:rule="!has(self.apiVersion) || (has(self.kind) && (self.apiVersion + '/' + self.kind) in ['user.dummy.example.org/v1alpha1/Human', 'user.dummy.example.org/v1alpha1/Machine'])",message="apiVersion and kind must be one of: user.dummy.example.org/v1alpha1 Human, user.dummy.example.org/v1alpha1 Machine"
// +kubebuilder:validation:Optional
`,
				},
			},
		},
		"MultiKindAmbiguousKind_namespaced": {
			args: args{
				crdScope: CRDScopeNamespaced,
				t:        types.NewTypeName(token.NoPos, tp, "Params", types.Universe.Lookup("string").Type()),
				f: &Field{
					Name: name.NewFromCamel("GrantID"),
					Reference: &config.Reference{
						Type:       "github.com/upbound/official-providers/provider-dummy/apis/project/v1alpha1.Grant",
						APIVersion: "project.dummy.example.org/v1alpha1",
						AdditionalTargets: []config.ReferenceTarget{{
							Type:       "github.com/upbound/official-providers/provider-dummy/apis/user/v1alpha1.Grant",
							APIVersion: "user.dummy.example.org/v1alpha1",
						}, {
							Type:       "github.com/upbound/official-providers/provider-dummy/apis/user/v1alpha1.Machine",
							APIVersion: "user.dummy.example.org/v1alpha1",
						}},
					},
					FieldType: types.NewPointer(types.Universe.Lookup("string").Type()),
				},
			}, want: want{
				outFields: []*types.Var{
					types.NewField(token.NoPos, tp, "GrantIDRef", types.NewPointer(typeNamespacedKindReferenceField), false),
					types.NewField(token.NoPos, tp, "GrantIDSelector", types.NewPointer(typeNamespacedKindSelectorField), false),
				},
				outTags: []string{
					`json:"grantIdRef,omitempty" tf:"-"`,
					`json:"grantIdSelector,omitempty" tf:"-"`,
				},
				outComments: twtypes.Comments{
					"github.com/crossplane/upjet/v2/pkg/types.Params:GrantIDRef": `// Reference to a Grant in project (default), Grant in user or Machine in user to populate grantId.
// Set kind, and apiVersion if kind alone is ambiguous, to choose a target other than the default.
// +kubebuilder:validation:XValidation:rule="!has(self.kind) || self.kind in ['Grant', 'Machine']",message="kind must be one of Grant, Machine"
// +kubebuilder:validation:XValidation:rule="!has(self.apiVersion) || (has(self.kind) && (self.apiVersion + '/' + self.kind) in ['project.dummy.example.org/v1alpha1/Grant', 'user.dummy.example.org/v1alpha1/Grant', 'user.dummy.example.org/v1alpha1/Machine'])",message="apiVersion and kind must be one of: project.dummy.example.org/v1alpha1 Grant, user.dummy.example.org/v1alpha1 Grant, user.dummy.example.org/v1alpha1 Machine"
// +kubebuilder:validation:XValidation:rule="!has(self.kind) || has(self.apiVersion) || !(self.kind in ['Grant'])",message="apiVersion is required when kind is Grant"
// +kubebuilder:validation:Optional
`,
					"github.com/crossplane/upjet/v2/pkg/types.Params:GrantIDSelector": `// Selector for a Grant in project (default), Grant in user or Machine in user to populate grantId.
// Set kind, and apiVersion if kind alone is ambiguous, to choose a target other than the default.
// +kubebuilder:validation:XValidation:rule="!has(self.kind) || self.kind in ['Grant', 'Machine']",message="kind must be one of Grant, Machine"
// +kubebuilder:validation:XValidation:rule="!has(self.apiVersion) || (has(self.kind) && (self.apiVersion + '/' + self.kind) in ['project.dummy.example.org/v1alpha1/Grant', 'user.dummy.example.org/v1alpha1/Grant', 'user.dummy.example.org/v1alpha1/Machine'])",message="apiVersion and kind must be one of: project.dummy.example.org/v1alpha1 Grant, user.dummy.example.org/v1alpha1 Grant, user.dummy.example.org/v1alpha1 Machine"
// +kubebuilder:validation:XValidation:rule="!has(self.kind) || has(self.apiVersion) || !(self.kind in ['Grant'])",message="apiVersion is required when kind is Grant"
// +kubebuilder:validation:Optional
`,
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			g := &Builder{
				comments: twtypes.Comments{},
				scope:    tc.args.crdScope,
			}
			gotFields, gotTags := g.generateReferenceFields(tc.args.t, tc.args.f)
			if diff := cmp.Diff(tc.want.outFields, gotFields, cmp.Comparer(func(a, b *types.Var) bool {
				return a.String() == b.String()
			})); diff != "" {
				t.Errorf("generateReferenceFields(): fields: +got, -want: %s", diff)
			}
			if diff := cmp.Diff(tc.want.outTags, gotTags); diff != "" {
				t.Errorf("generateReferenceFields(): tags: +got, -want: %s", diff)
			}
			if diff := cmp.Diff(tc.want.outComments, g.comments); diff != "" {
				t.Errorf("generateReferenceFields(): comments: +got, -want: %s", diff)
			}
		})
	}
}
