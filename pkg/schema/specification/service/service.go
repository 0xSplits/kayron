package service

import (
	"github.com/0xSplits/kayron/pkg/schema/specification/service/deploy"
	"github.com/0xSplits/kayron/pkg/schema/specification/service/docker"
	"github.com/0xSplits/kayron/pkg/schema/specification/service/github"
	"github.com/xh3b4sd/tracer"
)

type Service struct {
	Docker docker.Docker `json:"docker,omitempty" yaml:"docker,omitempty"`
	GitHub github.Github `json:"github,omitempty" yaml:"github,omitempty"`
	Deploy deploy.Deploy `json:"deploy,omitempty" yaml:"deploy,omitempty"`
}

func (s Service) Empty() bool {
	return s.Deploy.Empty() && s.GitHub.Empty() && s.Deploy.Empty()
}

func (s Service) Verify() error {
	{
		err := s.Docker.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		err := s.GitHub.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		err := s.Deploy.Verify()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
