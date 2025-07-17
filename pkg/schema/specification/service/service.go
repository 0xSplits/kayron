package service

import (
	"github.com/0xSplits/kayron/pkg/schema/specification/service/deploy"
	"github.com/0xSplits/kayron/pkg/schema/specification/service/docker"
	"github.com/0xSplits/kayron/pkg/schema/specification/service/github"
	"github.com/xh3b4sd/tracer"
)

type Service struct {
	Docker docker.Docker `yaml:"docker,omitempty"`
	Github github.Github `yaml:"github,omitempty"`
	Deploy deploy.Deploy `yaml:"deploy,omitempty"`
}

func (s Service) Empty() bool {
	return s.Docker.Empty() && s.Github.Empty() && s.Deploy.Empty()
}

func (s Service) Verify() error {
	if s.Docker.Empty() {
		return tracer.Mask(serviceDockerEmptyError)
	} else {
		err := s.Docker.Verify()
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

	if s.Deploy.Empty() {
		// the default deployment strategy is to rollout the latest release
	} else {
		err := s.Deploy.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
