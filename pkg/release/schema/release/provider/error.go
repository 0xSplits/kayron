package provider

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

//
//
//

var providerNameError = &tracer.Error{
	Description: "The provider configuration requires the provider name to be \"cloudformation\".",
}

func IsProviderName(err error) bool {
	return errors.Is(err, providerNameError)
}
