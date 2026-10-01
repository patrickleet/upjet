// SPDX-FileCopyrightText: 2023 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package reference

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/registry"
	"github.com/crossplane/upjet/v2/pkg/types"
)

const (
	extractorPackagePath      = "github.com/crossplane/upjet/v2/pkg/resource"
	extractResourceIDFuncPath = extractorPackagePath + ".ExtractResourceID()"
	fmtExtractParamFuncPath   = extractorPackagePath + `.ExtractParamPath("%s",%t)`
)

// Injector resolves references using provider metadata
type Injector struct {
	ModulePath        string
	ProviderShortName string
	// RootGroup is the provider's root API group. It's used to set the
	// APIVersion of multi-kind reference targets.
	RootGroup string
}

// NewInjector initializes a new Injector
func NewInjector(apisModulePath string) *Injector {
	return &Injector{
		ModulePath: apisModulePath,
	}
}

func getExtractorFuncPath(r *config.Resource, sourceAttr string) string {
	switch sourceAttr {
	// value extractor from status.atProvider.id
	case "id":
		return extractResourceIDFuncPath
	// value extractor from spec.forProvider.<attr>
	default:
		for _, n := range r.ExternalName.OmittedFields {
			if sourceAttr == n {
				return ""
			}
		}
		s, ok := r.TerraformResource.Schema[sourceAttr]
		if !ok {
			return ""
		}
		return fmt.Sprintf(fmtExtractParamFuncPath, sourceAttr, types.IsObservation(s))
	}
}

// InjectReferences injects cross-resource references using the
// provider metadata scraped from the Terraform registry.
func (rr *Injector) InjectReferences(configResources map[string]*config.Resource) error { //nolint:gocyclo
	for n, r := range configResources {
		m := configResources[n].MetaResource
		if m == nil {
			continue
		}

		for i, re := range m.Examples {
			pm, err := paveExampleManifest(re.Manifest)
			if err != nil {
				return errors.Wrapf(err, "cannot pave example manifest for resource: %s", n)
			}
			resolutionContext, err := PrepareLocalResolutionContext(re, NewRefParts(n, re.Name).GetResourceName(false))
			if err != nil {
				return errors.Wrapf(err, "cannot prepare local resolution context for resource: %s", n)
			}
			if err := rr.ResolveReferencesOfPaved(pm, resolutionContext); err != nil {
				return errors.Wrapf(err, "cannot resolve references of resource with local examples context: %s", n)
			}
			if err := rr.storeResolvedDependencies(&m.Examples[i], resolutionContext.Context); err != nil {
				return errors.Wrapf(err, "cannot store resolved dependencies for resource: %s", n)
			}
			for targetAttr, ref := range re.References {
				// if a reference is already configured for the target attribute
				if _, ok := r.References[targetAttr]; ok {
					continue
				}
				parts := getRefParts(ref)
				// if nil or a references to a nested configuration block
				if parts == nil || strings.Contains(parts.Attribute, ".") || strings.Contains(parts.Attribute, "[") {
					continue
				}
				if _, ok := configResources[parts.Resource]; !ok {
					continue
				}
				r.References[targetAttr] = config.Reference{
					TerraformName: parts.Resource,
					Extractor:     getExtractorFuncPath(configResources[parts.Resource], parts.Attribute),
				}
			}
		}
	}
	return nil
}

func (rr *Injector) storeResolvedDependencies(re *registry.ResourceExample, context map[string]*PavedWithManifest) error {
	for rn, pm := range context {
		buff, err := pm.Paved.MarshalJSON()
		if err != nil {
			return errors.Wrapf(err, "cannot marshal paved as JSON: %s", rn)
		}
		if _, ok := re.Dependencies[rn]; ok {
			re.Dependencies[rn] = string(buff)
		}
	}
	return nil
}

func (rr *Injector) getTypePath(tfName string, configResources map[string]*config.Resource) (string, error) {
	r := configResources[tfName]
	if r == nil {
		return "", errors.Errorf("cannot find configuration for Terraform resource: %s", tfName)
	}
	shortGroup := r.ShortGroup
	if len(shortGroup) == 0 {
		shortGroup = rr.ProviderShortName
	}
	return fmt.Sprintf("%s/%s/%s.%s", rr.ModulePath, shortGroup, r.Version, r.Kind), nil
}

func (rr *Injector) getAPIVersion(tfName string, configResources map[string]*config.Resource) (string, error) {
	r := configResources[tfName]
	if r == nil {
		return "", errors.Errorf("cannot find configuration for Terraform resource: %s", tfName)
	}
	if rr.RootGroup == "" {
		return "", errors.Errorf("cannot determine the API version of Terraform resource %s: the root group is not set", tfName)
	}
	group := rr.RootGroup
	if r.ShortGroup != "" {
		group = strings.ToLower(r.ShortGroup) + "." + rr.RootGroup
	}
	return group + "/" + r.Version, nil
}

// setTargetTypes sets the Type and APIVersion of the targets of a multi-kind
// reference that are configured by TerraformName, and validates the targets.
func (rr *Injector) setTargetTypes(ref *config.Reference, configResources map[string]*config.Resource) error {
	if ref.TerraformName != "" && ref.APIVersion == "" {
		v, err := rr.getAPIVersion(ref.TerraformName, configResources)
		if err != nil {
			return err
		}
		ref.APIVersion = v
	}
	// copy the targets so that a configuration shared between
	// resources isn't modified in place.
	targets := make([]config.ReferenceTarget, len(*ref.AdditionalTargets))
	for i, t := range *ref.AdditionalTargets {
		var err error
		if targets[i], err = rr.setTargetType(t, configResources); err != nil {
			return err
		}
	}
	ref.AdditionalTargets = &targets
	return ref.ValidateTargets()
}

// setTargetType sets the Type and APIVersion of a target configured by
// TerraformName, unless they're already set.
func (rr *Injector) setTargetType(t config.ReferenceTarget, configResources map[string]*config.Resource) (config.ReferenceTarget, error) {
	if t.TerraformName == "" {
		return t, nil
	}
	var err error
	if t.Type == "" {
		if t.Type, err = rr.getTypePath(t.TerraformName, configResources); err != nil {
			return t, err
		}
	}
	if t.APIVersion == "" {
		if t.APIVersion, err = rr.getAPIVersion(t.TerraformName, configResources); err != nil {
			return t, err
		}
	}
	return t, nil
}

// SetReferenceTypes resolves reference types of configured references
// using their TerraformNames.
func (rr *Injector) SetReferenceTypes(configResources map[string]*config.Resource) error {
	for name, r := range configResources {
		for attr, ref := range r.References {
			if ref.Type == "" && ref.TerraformName != "" { //nolint:staticcheck // still handling deprecated field behavior
				crdTypePath, err := rr.getTypePath(ref.TerraformName, configResources)
				if err != nil {
					return errors.Wrap(err, "cannot set reference types")
				}
				// TODO(aru): if type mapper cannot provide a mapping,
				// currently we remove the reference. Once,
				// we have type mapper implementations available
				// for all providers, then we can keep the refs
				// instead of removing them, and expect resulting
				// compile errors to be fixed by making the types
				// available to the type mapper.
				if crdTypePath == "" {
					delete(r.References, attr)
					continue
				}
				ref.Type = crdTypePath //nolint:staticcheck // still handling deprecated field behavior
				r.References[attr] = ref
			}
			if len(ref.Targets()) > 0 {
				if err := rr.setTargetTypes(&ref, configResources); err != nil {
					return errors.Wrapf(err, "cannot set the reference targets of %s.%s", name, attr)
				}
				r.References[attr] = ref
			}
		}
	}
	return nil
}
