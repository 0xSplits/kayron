package deploy

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var deploymentStrategyError = &tracer.Error{
	Description: "The deployment configuration requires only one strategy to be provided.",
}

func IsDeploymentStrategy(err error) bool {
	return errors.Is(err, deploymentStrategyError)
}
