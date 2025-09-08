package awsauthn

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Aws is a standard config object containing valid AWS credentials of any
	// kind, e.g. access keys or SSO tokens.
	Aws aws.Config
}

// Keychain is a container registry authenticator using the standard AWS
// credentials without requiring the host file system. Most keychain
// implementations work with the host's local .docker/config.json, which is not
// useful if we already have a set of standard AWS credentials setup for our
// program. This keychain implementation works by requesting an authoriczation
// token from ECR for the container registry identified by the provided AWS
// credentials, which usually resolves to the AWS account within which those AWS
// credentials have been defined.
type Keychain struct {
	ecr *ecr.Client
}

func New(c Config) *Keychain {
	if c.Aws.Region == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Aws must not be empty", c)))
	}

	return &Keychain{
		ecr: ecr.NewFromConfig(c.Aws),
	}
}

func (k *Keychain) Resolve(r authn.Resource) (authn.Authenticator, error) {
	var err error

	// Fetch the standard authorization token from ECR without the need to
	// interact with any file system or host config files.

	var inp *ecr.GetAuthorizationTokenInput
	{
		inp = &ecr.GetAuthorizationTokenInput{}
	}

	var out *ecr.GetAuthorizationTokenOutput
	{
		out, err = k.ecr.GetAuthorizationToken(context.Background(), inp)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	// Guard against an invalid response from the ECR APIs.

	if len(out.AuthorizationData) != 1 || out.AuthorizationData[0].AuthorizationToken == nil {
		return nil, tracer.Mask(authorizationTokenError, tracer.Context{Key: "reason", Value: "no token"})
	}

	// Decode the received API response in the form base64("AWS:<password>") in
	// order to create a basic auth authenticator that the standard container
	// registry interface understands.

	var b64 string
	{
		b64 = *out.AuthorizationData[0].AuthorizationToken
	}

	var dec []byte
	{
		dec, err = base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var spl []string
	{
		spl = strings.SplitN(string(dec), ":", 2)
		if len(spl) != 2 {
			return nil, tracer.Mask(authorizationTokenError, tracer.Context{Key: "reason", Value: "wrong format"})
		}
	}

	// Return a basic auth authenticator for some standard container registry.

	return &authn.Basic{Username: spl[0], Password: spl[1]}, nil
}
