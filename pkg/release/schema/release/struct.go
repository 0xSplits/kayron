package release

import (
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/release/docker"
	"github.com/0xSplits/kayron/pkg/release/schema/release/github"
	"github.com/0xSplits/kayron/pkg/release/schema/release/labels"
	"github.com/0xSplits/kayron/pkg/release/schema/release/provider"
	"github.com/xh3b4sd/tracer"
)

type Struct struct {
	Deploy   deploy.Struct   `yaml:"deploy,omitempty"`
	Docker   docker.String   `yaml:"docker,omitempty"`
	Github   github.String   `yaml:"github,omitempty"`
	Labels   labels.Struct   `yaml:"-"`
	Provider provider.String `yaml:"provider,omitempty"`
}

func (s Struct) Empty() bool {
	return s.Docker.Empty() && s.Github.Empty() && s.Deploy.Empty() && s.Labels.Empty() && s.Provider.Empty()
}

func (s Struct) Verify() error {
	err := s.verify()
	if err != nil {
		return tracer.Mask(err,
			tracer.Context{Key: "index", Value: s.Labels.Block},
			tracer.Context{Key: "file", Value: s.Labels.Source},
		)
	}

	return nil
}

func (s Struct) verify() error {
	if s.Deploy.Empty() {
		return tracer.Mask(releaseDeployEmptyError)
	} else {
		err := s.Deploy.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if s.Github.Empty() {
		return tracer.Mask(releaseGithubEmptyError)
	} else {
		err := s.Github.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if s.Labels.Empty() {
		return tracer.Mask(releaseLabelsEmptyError)
	} else {
		err := s.Labels.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if s.Docker.Empty() && s.Provider.Empty() {
		return tracer.Mask(releaseProviderEmptyError)
	} else {
		if !s.Docker.Empty() {
			err := s.Docker.Verify()
			if err != nil {
				return tracer.Mask(err)
			}
		}

		if !s.Provider.Empty() {
			err := s.Provider.Verify()
			if err != nil {
				return tracer.Mask(err)
			}
		}
	}

	return nil
}
