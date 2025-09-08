package awsauthn

import (
	"github.com/xh3b4sd/tracer"
)

var authorizationTokenError = &tracer.Error{
	Description: "This critical error indicates that no authorization token could be obtained by the resolver, which means that the keychain does not know how to proceed safely.",
}
