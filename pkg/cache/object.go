package cache

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// kind is a private type for private object properties in order to guarantee
// differentiated access patterns for infrastructure and service releases.
type kind string

const (
	Infrastructure kind = "infrastructure"
	Service        kind = "service"
)

// Object combines the associated artifact and release information in one
// addressable cache object. The artifact information are written and the
// release information are read throughout the reconciliation loops.
type Object struct {
	Artifact artifact.Struct
	Release  release.Struct

	ind int
	kin kind
}

func (o Object) Name() string {
	if o.kin == Infrastructure {
		return o.Release.Github.String()
	}

	if o.kin == Service {
		return o.Release.Docker.String()
	}

	return ""
}

// Parameter returns the CloudFormation stack parameter key for this release
// artifact. The parameter keys generated here have to be supported in the
// CloudFormation template being deployed, e.g. InfrastructureVersion,
// SpectaVersion.
func (o Object) Parameter() string {
	return fmt.Sprintf("%sVersion", cases.Title(language.English).String(o.Name()))
}
