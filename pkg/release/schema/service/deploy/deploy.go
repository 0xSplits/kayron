package deploy

import (
	"github.com/0xSplits/kayron/pkg/release/schema/service/deploy/branch"
	"github.com/0xSplits/kayron/pkg/release/schema/service/deploy/release"
	"github.com/0xSplits/kayron/pkg/release/schema/service/deploy/suspend"
	"github.com/0xSplits/kayron/pkg/release/schema/service/deploy/webhook"
	"github.com/xh3b4sd/tracer"
)

// Struct defines exactly one mutually exclusive declaration of either Branch,
// Release, Suspend or Webhook as required deployment instruction.
type Struct struct {
	Branch  branch.String  `yaml:"branch,omitempty"`
	Release release.String `yaml:"release,omitempty"`
	Suspend suspend.Bool   `yaml:"suspend,omitempty"`
	Webhook webhook.Slice  `yaml:"webhook,omitempty"`
}

func (s Struct) Empty() bool {
	return s.Branch.Empty() && s.Release.Empty() && s.Suspend.Empty() && s.Webhook.Empty()
}

func (s Struct) String() string {
	if !s.Branch.Empty() {
		return s.Branch.String()
	}

	if !s.Release.Empty() {
		return s.Release.String()
	}

	if !s.Suspend.Empty() {
		return s.Suspend.String()
	}

	if !s.Webhook.Empty() {
		return s.Webhook.String()
	}

	return ""
}

func (s Struct) Verify() error {
	// Reject deployment configurations that define more than one strategy.
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
