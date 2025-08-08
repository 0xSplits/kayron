package infrastructure

import (
	"bytes"
	"context"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/xh3b4sd/tracer"
)

func (i *Infrastructure) putObj(pat string, byt []byte) error {
	var err error

	var key string
	{
		key, err = i.envKey(pat)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		i.log.Log(
			"level", "debug",
			"message", "uploading cloudformation template",
			"bucket", Bucket,
			"key", key,
		)
	}

	// Make sure we respect the dry run flag when attempting to upload any
	// template to S3, because "dry run" effectively means "read only". So
	// if the dry run flag is set in e.g. the operator's integration test,
	// then we want to emit the logs, but prevent making the network calls.

	var inp *s3.PutObjectInput
	{
		inp = &s3.PutObjectInput{
			Bucket:      aws.String(Bucket),
			Key:         aws.String(key),
			Body:        bytes.NewReader(byt),
			ContentType: aws.String("application/x-yaml"),
		}
	}

	if !i.dry {
		_, err = i.as3.PutObject(context.Background(), inp)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (i *Infrastructure) envKey(pat string) (string, error) {
	key, err := filepath.Rel(Directory, pat)
	if err != nil {
		return "", tracer.Mask(err)
	}

	return filepath.Join(i.env, key), nil
}
