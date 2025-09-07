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
	Aws aws.Config
}

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

	if len(out.AuthorizationData) != 1 || out.AuthorizationData[0].AuthorizationToken == nil {
		return nil, tracer.Mask(authorizationTokenError, tracer.Context{Key: "reason", Value: "no token"})
	}

	var b64 string
	{
		b64 = *out.AuthorizationData[0].AuthorizationToken // base64("AWS:<password>")
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

	return &authn.Basic{Username: spl[0], Password: spl[1]}, nil
}
