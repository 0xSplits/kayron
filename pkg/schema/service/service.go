package service

import (
	"github.com/0xSplits/kayron/pkg/schema/service/deploy"
	"github.com/0xSplits/kayron/pkg/schema/service/docker"
	"github.com/0xSplits/kayron/pkg/schema/service/github"
)

// Service describes one independently deployable unit, defined by a Docker
// image and a GitHub repository.
type Service struct {
	Docker docker.Docker `json:"docker,omitempty" yaml:"docker,omitempty"`
	GitHub github.Github `json:"github,omitempty" yaml:"github,omitempty"`
	Deploy deploy.Deploy `json:"deploy,omitempty" yaml:"deploy,omitempty"`
}
