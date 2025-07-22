package service

import (
	"github.com/0xSplits/kayron/pkg/release/schema/service/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/service/docker"
	"github.com/0xSplits/kayron/pkg/release/schema/service/github"
	"github.com/0xSplits/kayron/pkg/release/schema/service/labels"
	"github.com/0xSplits/kayron/pkg/release/schema/service/provider"
	"github.com/xh3b4sd/tracer"
)

type Service struct {
	Deploy   deploy.Struct   `yaml:"deploy,omitempty"`
	Docker   docker.String   `yaml:"docker,omitempty"`
	Github   github.String   `yaml:"github,omitempty"`
	Labels   labels.Struct   `yaml:"-"`
	Provider provider.String `yaml:"provider,omitempty"`
}

func (s Service) Empty() bool {
	return s.Docker.Empty() && s.Github.Empty() && s.Deploy.Empty() && s.Labels.Empty() && s.Provider.Empty()
}

func (s Service) Verify() error {
	err := s.verify()
	if err != nil {
		return tracer.Mask(err, tracer.Context{Key: "file", Value: s.Labels.Source})
	}

	return nil
}

func (s Service) verify() error {
	if s.Deploy.Empty() {
		return tracer.Mask(serviceDeployEmptyError)
	} else {
		err := s.Deploy.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if s.Github.Empty() {
		return tracer.Mask(serviceGithubEmptyError)
	} else {
		err := s.Github.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if s.Labels.Empty() {
		return tracer.Mask(serviceLabelsEmptyError)
	} else {
		err := s.Labels.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if s.Docker.Empty() && s.Provider.Empty() {
		return tracer.Mask(serviceProviderEmptyError)
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
