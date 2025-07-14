package deploy

import (
	"github.com/0xSplits/kayron/pkg/schema/service/deploy/branch"
	"github.com/0xSplits/kayron/pkg/schema/service/deploy/release"
	"github.com/0xSplits/kayron/pkg/schema/service/deploy/suspend"
	"github.com/0xSplits/kayron/pkg/schema/service/deploy/webhook"
)

// Deploy defines exactly one mutually exclusive declaration of either Branch,
// Release, Suspend or Webhook as required deployment instruction.
type Deploy struct {
	// Branch triggers branch-based deployment for test environments. For
	// instance, it is possible to instruct Kayron to deploy Specta according to
	// all changes of the specified feature branch in the Specta repository, and
	// deploy those changes to the test environment matching the current branch.
	// Note that this overwrite is only considered valid for existing test
	// environments.
	Branch branch.Branch `yaml:"branch,omitempty"`

	// Release must be a semver version string representing the respective Github
	// release tag. The tag format is required to contain a leading "v" prefix and
	// optional "-" separator for providing additional metadata.
	//
	//     v0.1.0            the very first development release for new projects
	//     v1.8.2            the fully qualified first major release for stable APIs
	//     v1.8.3-ffce1e2    the metadata version for third party projects like Alloy
	//
	Release release.Release `yaml:"release,omitempty"`

	// Suspend disables any further reconciliation of this service indefinitely.
	Suspend suspend.Suspend `yaml:"suspend,omitempty"`

	// Webhook contains a list of alternative deployment mechanisms. Each webhook
	// provided here is invoked to deploy e.g. our frontends in Vercel. The format
	// of those webhook definitions requires the usage of a prefix for the HTTP
	// method that this webhook should be called with. It is further required to
	// provide a HTTPs URL. Failed webhook calls may be retried and eventually be
	// reported as terminal failure.
	//
	//     POST:https://{{DNS}}/{{PATH}}
	//
	Webhook webhook.Webhook `yaml:"webhook,omitempty"`
}

func (d Deploy) Empty() bool {
	return d.Branch.Empty() && d.Release.Empty() && d.Suspend.Empty() && d.Webhook.Empty()
}

func (d Deploy) Verify() error {
	return nil
}
