package deploy

import (
	"github.com/0xSplits/kayron/pkg/schema/specification/service/deploy/branch"
	"github.com/0xSplits/kayron/pkg/schema/specification/service/deploy/release"
	"github.com/0xSplits/kayron/pkg/schema/specification/service/deploy/suspend"
	"github.com/0xSplits/kayron/pkg/schema/specification/service/deploy/webhook"
	"github.com/xh3b4sd/tracer"
)

// Deploy defines exactly one mutually exclusive declaration of either Branch,
// Release, Suspend or Webhook as required deployment instruction.
type Deploy struct {
	Branch  branch.Branch    `yaml:"branch,omitempty"`
	Release release.Release  `yaml:"release,omitempty"`
	Suspend suspend.Suspend  `yaml:"suspend,omitempty"`
	Webhook webhook.Webhooks `yaml:"webhook,omitempty"`
}

func (d Deploy) Empty() bool {
	return d.Branch.Empty() && d.Release.Empty() && d.Suspend.Empty() && d.Webhook.Empty()
}

func (d Deploy) Verify() error {
	// Reject deployment configurations that define more than one strategy.
	{
		lis := enabled(d.Branch, d.Release, d.Suspend, d.Webhook)
		if len(lis) > 1 {
			return tracer.Mask(deploymentStrategyError, tracer.Context{Key: "enabled", Value: lis})
		}
	}

	if !d.Branch.Empty() {
		err := d.Branch.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if !d.Release.Empty() {
		err := d.Release.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if !d.Suspend.Empty() {
		err := d.Suspend.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if !d.Webhook.Empty() {
		err := d.Webhook.Verify()
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
