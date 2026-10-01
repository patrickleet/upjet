// SPDX-FileCopyrightText: 2023 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package markers

import (
	"fmt"

	"github.com/crossplane/upjet/v2/pkg/config"
)

const (
	markerPrefixCrossplane = "+crossplane:"
)

var (
	markerPrefixRefType         = fmt.Sprintf("%sgenerate:reference:type=", markerPrefixCrossplane)
	markerPrefixRefExtractor    = fmt.Sprintf("%sgenerate:reference:extractor=", markerPrefixCrossplane)
	markerPrefixRefFieldName    = fmt.Sprintf("%sgenerate:reference:refFieldName=", markerPrefixCrossplane)
	markerPrefixRefSelectorName = fmt.Sprintf("%sgenerate:reference:selectorFieldName=", markerPrefixCrossplane)
	markerPrefixRefAPIVersion   = fmt.Sprintf("%sgenerate:reference:apiVersion=", markerPrefixCrossplane)
)

// defaultExtractor is the extractor angryjet uses when none is configured.
// Multi-kind references spell it out when another target has an extractor,
// because their extractor markers must be either absent or one per target.
const defaultExtractor = "github.com/crossplane/crossplane-runtime/v2/pkg/reference.ExternalName()"

// CrossplaneOptions represents the Crossplane marker options that upjet
// would need to interact
type CrossplaneOptions struct {
	config.Reference
}

func (o CrossplaneOptions) String() string {
	m := ""

	if targets := o.Targets(); len(targets) > 0 {
		m += multiKindMarkers(targets)
	} else if o.Type != "" {
		m += fmt.Sprintf("%s%s\n", markerPrefixRefType, o.Type)
	}
	if o.Extractor != "" && len(o.AdditionalTargets) == 0 {
		m += fmt.Sprintf("%s%s\n", markerPrefixRefExtractor, o.Extractor)
	}
	if o.RefFieldName != "" {
		m += fmt.Sprintf("%s%s\n", markerPrefixRefFieldName, o.RefFieldName)
	}
	if o.SelectorFieldName != "" {
		m += fmt.Sprintf("%s%s\n", markerPrefixRefSelectorName, o.SelectorFieldName)
	}

	return m
}

// multiKindMarkers returns the type, apiVersion and (if any target has one)
// extractor markers of a multi-kind reference, one group per target. The
// first target is the default.
func multiKindMarkers(targets []config.ReferenceTarget) string {
	withExtractors := false
	for _, t := range targets {
		withExtractors = withExtractors || t.Extractor != ""
	}
	m := ""
	for _, t := range targets {
		m += fmt.Sprintf("%s%s\n", markerPrefixRefType, t.Type)
		m += fmt.Sprintf("%s%s\n", markerPrefixRefAPIVersion, t.APIVersion)
		if !withExtractors {
			continue
		}
		e := t.Extractor
		if e == "" {
			e = defaultExtractor
		}
		m += fmt.Sprintf("%s%s\n", markerPrefixRefExtractor, e)
	}
	return m
}
