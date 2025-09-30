package cache

import (
	"fmt"
	"strings"

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

// Domain returns the hash based testing domain for preview deployments, or an
// empty string for any other main release and non-testing environment.
//
//	1d0fd508.lite.testing.splits.org
func (o Object) Domain(env string) string {
	// Note that we filter the domain name creation by preview deployments,
	// because at the time of writing we do not have any convenient way to tell
	// whether this release artifact is exposed to the internet via DNS. Right now
	// we only know that for certain in case of preview deployments, because their
	// sole purpose is to be exposed to the internet.

	if !o.Preview() {
		return ""
	}

	return fmt.Sprintf("%s.%s.%s.splits.org",
		o.Release.Labels.Hash.Lower(),

		// Note that this is a dirty hack to make preview deployments work today for
		// existing services that already work using certain incosnistencies between
		// repository and domain names. E.g. we have "splits-lite" in Github, but
		// use just "lite.testing.splits.org". A better way of doing this would be
		// to allow for some kind of domain configuration in the release definition,
		// so that we can remove this magical string replacement below.
		strings.TrimPrefix(o.Release.Docker.String(), "splits-"),

		env,
	)
}

// Drift tries to detect a single valid state drift, in order to allow allow the
// operator chain to execute. Our current policy requires the following
// conditions to be true for a valid state drift using the ready flag. If the
// ready flag is false, then only the first three conditions are being
// evaluated. That way we can differentiate between ready and waiting state
// drift. Note that the set of waiting release artifacts may also contain ready
// release artifacts.
//
//  1. the desired deployment must not be suspended
//
//  2. the current and desired state must not be equal
//
//  3. the desired state must not be empty
//
//  4. the container image for the desired state must be pushed
func (o Object) Drift(rea bool) bool {
	return !bool(o.Release.Deploy.Suspend) &&
		o.Artifact.Drift() &&
		!o.Artifact.Empty() &&
		(!rea || o.Artifact.Valid())
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
	return fmt.Sprintf("%sVersion", strings.Map(mapFnc, cases.Title(language.English).String(o.Name())))
}

// Preview returns whether this release artifact is for an automatically
// injected preview deployment. The determining factor here is whether a preview
// hash exists.
func (o Object) Preview() bool {
	return !o.Release.Labels.Hash.Empty()
}

// Version returns the desired state of this artifact's release version if the
// following two conditions are true. Failed artifact conditions and deployment
// suspensions will then yield the artifact version of the scheduler's current
// state.
//
//  1. the artifact condition must be true
//
//  2. the deployment strategy must not be suspended
//
// E.g. we may observe the new version of a service v0.5.0, for which there does
// no docker image exist yet. In this case we defer to the old version v0.4.0,
// which is currently deployed. Once the docker image for v0.5.0 has been
// confirmed inside of the underlying container registry, we will yield the new
// desired version.
func (o Object) Version() string {
	if bool(o.Release.Deploy.Suspend) || !o.Artifact.Valid() {
		return o.Artifact.Scheduler.Current
	}

	return o.Artifact.Reference.Desired
}

// mapFnc signals the removal of dashes and spaces only.
func mapFnc(r rune) rune {
	if r == '-' || r == ' ' {
		return -1
	}

	return r
}
