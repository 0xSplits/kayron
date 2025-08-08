package deploy

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var deploymentStrategyError = &tracer.Error{
	Description: "The deployment configuration requires only one strategy to be provided.",
	Context: []tracer.Context{
		{Key: "reason", Value: "more than one deployment strategy was found to be enabled"},
	},
}

func IsDeploymentStrategy(err error) bool {
	return errors.Is(err, deploymentStrategyError)
}
