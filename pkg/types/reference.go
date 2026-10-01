// SPDX-FileCopyrightText: 2023 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"go/token"
	"go/types"
	"reflect"
	"strings"

	"k8s.io/utils/ptr"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/types/comments"
	"github.com/crossplane/upjet/v2/pkg/types/markers"
	"github.com/crossplane/upjet/v2/pkg/types/markers/kubebuilder"
	"github.com/crossplane/upjet/v2/pkg/types/name"
)

const (
	// PackagePathXPCommonAPIs is the go path for the Crossplane Runtime package
	// with common APIs
	PackagePathXPCommonAPIs = "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	// PackagePathXPV2CommonAPIs is the go path for the Crossplane Runtime package
	// with common v2 APIs
	PackagePathXPV2CommonAPIs = "github.com/crossplane/crossplane-runtime/v2/apis/common/v2"
	// PackagePathKindRef is the go path for the package with the reference
	// and selector types of multi-kind references.
	PackagePathKindRef = "github.com/crossplane/upjet/v2/pkg/resource/kindref"
)

// Types to use from by reference generator.
var (
	typeReferenceField types.Type = types.NewNamed(
		types.NewTypeName(token.NoPos, types.NewPackage(PackagePathXPCommonAPIs, "v1"), "Reference", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	typeSelectorField types.Type = types.NewNamed(
		types.NewTypeName(token.NoPos, types.NewPackage(PackagePathXPCommonAPIs, "v1"), "Selector", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	typeSecretKeySelector types.Type = types.NewNamed(
		types.NewTypeName(token.NoPos, types.NewPackage(PackagePathXPCommonAPIs, "v1"), "SecretKeySelector", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	typeSecretReference types.Type = types.NewNamed(
		types.NewTypeName(token.NoPos, types.NewPackage(PackagePathXPCommonAPIs, "v1"), "SecretReference", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	typeLocalSecretReference types.Type = types.NewNamed(
		types.NewTypeName(token.NoPos, types.NewPackage(PackagePathXPCommonAPIs, "v1"), "LocalSecretReference", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	typeLocalSecretKeySelector types.Type = types.NewNamed(
		types.NewTypeName(token.NoPos, types.NewPackage(PackagePathXPCommonAPIs, "v1"), "LocalSecretKeySelector", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	typeNamespacedReferenceField types.Type = types.NewNamed(
		types.NewTypeName(token.NoPos, types.NewPackage(PackagePathXPCommonAPIs, "v1"), "NamespacedReference", nil),
		types.NewStruct(nil, nil),
		nil,
	)

	typeNamespacedSelectorField types.Type = types.NewNamed(
		types.NewTypeName(token.NoPos, types.NewPackage(PackagePathXPCommonAPIs, "v1"), "NamespacedSelector", nil),
		types.NewStruct(nil, nil),
		nil,
	)
	typeKindReferenceField           types.Type = newKindRefType("Reference")
	typeKindSelectorField            types.Type = newKindRefType("Selector")
	typeNamespacedKindReferenceField types.Type = newKindRefType("NamespacedReference")
	typeNamespacedKindSelectorField  types.Type = newKindRefType("NamespacedSelector")

	commentOptional = &comments.Comment{
		Options: markers.Options{
			KubebuilderOptions: kubebuilder.Options{
				Required: ptr.To(false),
			},
		},
	}
)

func newKindRefType(name string) types.Type {
	return types.NewNamed(
		types.NewTypeName(token.NoPos, types.NewPackage(PackagePathKindRef, "kindref"), name, nil),
		types.NewStruct(nil, nil),
		nil,
	)
}

func (g *Builder) generateReferenceFields(t *types.TypeName, f *Field) (fields []*types.Var, tags []string) {
	_, isSlice := f.FieldType.(*types.Slice)

	rfn := name.ReferenceFieldName(f.Name, isSlice, f.Reference.RefFieldName)
	sfn := name.SelectorFieldName(f.Name, f.Reference.SelectorFieldName)

	refTag := fmt.Sprintf(`json:"%s,omitempty" tf:"-"`, rfn.LowerCamelComputed)
	selTag := fmt.Sprintf(`json:"%s,omitempty" tf:"-"`, sfn.LowerCamelComputed)

	if targets := f.Reference.Targets(); len(targets) > 0 {
		return g.generateMultiKindReferenceFields(t, f, targets, rfn, sfn), []string{refTag, selTag}
	}

	var tr types.Type
	if g.scope == CRDScopeCluster {
		tr = types.NewPointer(typeReferenceField)
	} else {
		tr = types.NewPointer(typeNamespacedReferenceField)
	}
	refComment := fmt.Sprintf("// Reference to a %s to populate %s.\n%s",
		friendlyTypeDescription(f.Reference.Type), f.Name.LowerCamelComputed, commentOptional.Build())
	selComment := fmt.Sprintf("// Selector for a %s to populate %s.\n%s",
		friendlyTypeDescription(f.Reference.Type), f.Name.LowerCamelComputed, commentOptional.Build())
	if isSlice {
		tr = types.NewSlice(typeNamespacedReferenceField)
		if g.scope == CRDScopeCluster {
			tr = types.NewSlice(typeReferenceField)
		}
		refComment = fmt.Sprintf("// References to %s to populate %s.\n%s",
			friendlyTypeDescription(f.Reference.Type), f.Name.LowerCamelComputed, commentOptional.Build())
		selComment = fmt.Sprintf("// Selector for a list of %s to populate %s.\n%s",
			friendlyTypeDescription(f.Reference.Type), f.Name.LowerCamelComputed, commentOptional.Build())
	}
	ref := types.NewField(token.NoPos, g.Package, rfn.Camel, tr, false)
	tsel := types.NewPointer(typeNamespacedSelectorField)
	if g.scope == CRDScopeCluster {
		tsel = types.NewPointer(typeSelectorField)
	}
	sel := types.NewField(token.NoPos, g.Package, sfn.Camel, tsel, false)

	g.comments.AddFieldComment(t, rfn.Camel, refComment)
	g.comments.AddFieldComment(t, sfn.Camel, selComment)
	f.TransformedName = rfn.LowerCamelComputed
	f.SelectorName = sfn.LowerCamelComputed

	return []*types.Var{ref, sel}, []string{refTag, selTag}
}

// generateMultiKindReferenceFields generates the reference and selector
// fields of a reference with more than one target. Their types carry an
// optional apiVersion and kind that choose the target, which the generated
// validation rules restrict to the configured targets.
func (g *Builder) generateMultiKindReferenceFields(t *types.TypeName, f *Field, targets []config.ReferenceTarget, rfn, sfn name.Name) []*types.Var {
	tr, tsel := typeKindReferenceField, typeKindSelectorField
	if g.scope != CRDScopeCluster {
		tr, tsel = typeNamespacedKindReferenceField, typeNamespacedKindSelectorField
	}
	descs := make([]string, len(targets))
	for i, tg := range targets {
		descs[i] = friendlyTypeDescription(tg.Type)
	}
	descs[0] += " (default)"
	desc := descs[0]
	if len(descs) > 1 {
		desc = strings.Join(descs[:len(descs)-1], ", ") + " or " + descs[len(descs)-1]
	}
	usage := "// Set kind, and apiVersion if kind alone is ambiguous, to choose a target other than the default.\n"
	rules := multiKindValidationRules(targets)
	refComment := fmt.Sprintf("// Reference to a %s to populate %s.\n%s%s%s",
		desc, f.Name.LowerCamelComputed, usage, rules, commentOptional.Build())
	selComment := fmt.Sprintf("// Selector for a %s to populate %s.\n%s%s%s",
		desc, f.Name.LowerCamelComputed, usage, rules, commentOptional.Build())

	ref := types.NewField(token.NoPos, g.Package, rfn.Camel, types.NewPointer(tr), false)
	sel := types.NewField(token.NoPos, g.Package, sfn.Camel, types.NewPointer(tsel), false)
	g.comments.AddFieldComment(t, rfn.Camel, refComment)
	g.comments.AddFieldComment(t, sfn.Camel, selComment)
	f.TransformedName = rfn.LowerCamelComputed
	f.SelectorName = sfn.LowerCamelComputed
	return []*types.Var{ref, sel}
}

// multiKindValidationRules returns the CEL validation markers that restrict
// the apiVersion and kind of a multi-kind reference or selector to its
// configured targets, and require an apiVersion for a kind that more than
// one target has.
func multiKindValidationRules(targets []config.ReferenceTarget) string {
	var kinds, pairs, ambiguous, pairDescs []string
	seen := map[string]bool{}
	for _, t := range targets {
		pairs = append(pairs, fmt.Sprintf("'%s/%s'", t.APIVersion, t.Kind()))
		pairDescs = append(pairDescs, t.APIVersion+" "+t.Kind())
		if seen[t.Kind()] {
			continue
		}
		seen[t.Kind()] = true
		kinds = append(kinds, t.Kind())
		if config.IsAmbiguousKind(targets, t.Kind()) {
			ambiguous = append(ambiguous, t.Kind())
		}
	}
	quote := func(ss []string) string {
		q := make([]string, len(ss))
		for i, s := range ss {
			q[i] = "'" + s + "'"
		}
		return strings.Join(q, ", ")
	}
	rules := fmt.Sprintf("// +kubebuilder:validation:XValidation:rule=\"!has(self.kind) || self.kind in [%s]\",message=\"kind must be one of %s\"\n",
		quote(kinds), strings.Join(kinds, ", "))
	rules += fmt.Sprintf("// +kubebuilder:validation:XValidation:rule=\"!has(self.apiVersion) || (has(self.kind) && (self.apiVersion + '/' + self.kind) in [%s])\",message=\"apiVersion and kind must be one of: %s\"\n",
		strings.Join(pairs, ", "), strings.Join(pairDescs, ", "))
	if len(ambiguous) > 0 {
		rules += fmt.Sprintf("// +kubebuilder:validation:XValidation:rule=\"!has(self.kind) || has(self.apiVersion) || !(self.kind in [%s])\",message=\"apiVersion is required when kind is %s\"\n",
			quote(ambiguous), strings.Join(ambiguous, " or "))
	}
	return rules
}

// TypePath returns go package path for the input type. This is a helper
// function to be used whenever this information is needed, like configuring to
// reference to a type. Should not be used if the type is in the same package as
// the caller.
func TypePath(i any) string {
	return reflect.TypeOf(i).PkgPath() + "." + reflect.TypeOf(i).Name()
}

func friendlyTypeDescription(path string) string {
	if !strings.Contains(path, ".") {
		return path
	}
	typeName := path[strings.LastIndex(path, ".")+1:]
	dirs := strings.Split(path, "/")
	groupName := dirs[len(dirs)-2]
	return fmt.Sprintf("%s in %s", typeName, groupName)
}
