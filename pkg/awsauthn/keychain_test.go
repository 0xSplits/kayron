package awsauthn

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
)

// Test_AwsAuthn_Interface ensures that this AWS Keychain implementation
// complies with the authn keychain interface. This test already fails at
// compile time if Keychain does not implement authn.Keychain.
func Test_AwsAuthn_Interface(t *testing.T) {
	var _ authn.Keychain = &Keychain{}
}
