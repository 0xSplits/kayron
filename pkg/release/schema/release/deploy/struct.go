package deploy

import (
	"fmt"

	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/branch"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/preview"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/suspend"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/webhook"
	"github.com/xh3b4sd/tracer"
)

// Struct defines exactly one mutually exclusive declaration of either Branch,
// Release, Suspend or Webhook as required deployment instruction. Struct may
// also define Preview as a testing environment only deployment mechanism for
// pull requests, which does not influence the main deployment strategy
// configuration mentioned above.
type Struct struct {
	Branch  branch.String  `yaml:"branch,omitempty"`
	Preview preview.Bool   `yaml:"preview,omitempty"`
	Release release.String `yaml:"release,omitempty"`
	Suspend suspend.Bool   `yaml:"suspend,omitempty"`
	Webhook webhook.Slice  `yaml:"webhook,omitempty"`
}

func (s Struct) Empty() bool {
	return s.Branch.Empty() && s.Preview.Empty() && s.Release.Empty() && s.Suspend.Empty() && s.Webhook.Empty()
}

func (s Struct) String() string {
	if !s.Branch.Empty() {
		return fmt.Sprintf("branch=%s", s.Branch.String())
	}

	// Note that Struct.Preview is not a deployment strategy in all environments,
	// so Struct.Preview does not contribute to the name/string representation of
	// this deployment strategy.

	if !s.Release.Empty() {
		return fmt.Sprintf("release=%s", s.Release.String())
	}

	if !s.Suspend.Empty() {
		return fmt.Sprintf("suspend=%s", s.Suspend.String())
	}

	if !s.Webhook.Empty() {
		return fmt.Sprintf("webhook=%s", s.Webhook.String())
	}

	return ""
}

func (s Struct) Verify() error {
	// Reject deployment configurations that define more than one strategy. Note
	// that s.Preview is not a deployment strategy and is therefore not considered
	// for this check.
	{
		lis := enabled(s.Branch, s.Release, s.Suspend, s.Webhook)
		if len(lis) > 1 {
			return tracer.Mask(deploymentStrategyError, tracer.Context{Key: "enabled", Value: lis})
		}
	}

	if !s.Branch.Empty() {
		err := s.Branch.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if !s.Preview.Empty() {
		err := s.Preview.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if !s.Release.Empty() {
		err := s.Release.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if !s.Suspend.Empty() {
		err := s.Suspend.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if !s.Webhook.Empty() {
		err := s.Webhook.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

// enabled returns a list of the names of the enabled deployment strategies.
func enabled(v ...Interface) []string {
	var lis []string

	for _, x := range v {
		if !x.Empty() {
			lis = append(lis, name(x))
		}
	}

	return lis
}
