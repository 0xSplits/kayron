package infrastructure

import (
	"bytes"
	"context"
	"io/fs"
	"path/filepath"

	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/0xSplits/kayron/pkg/roghfs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"
)

func (i *Infrastructure) Ensure() error {
	// Find the cached Github reference for the infrastructure repository
	// containing our CloudFormation templates.

	var ref string
	var ser service.Service
	for j := range i.ser.Length() {
		{
			ser, _ = i.ser.Search(j)
		}

		if ser.Github.String() != Repository {
			continue
		}

		{
			ref = artifact.ReferenceDesired(j)
		}
	}

	// Setup the Github file system implementation so that we can walk the
	// contents of the remote infrastructure repository.

	var gfs *roghfs.Roghfs
	{
		gfs = roghfs.New(roghfs.Config{
			Bas: afero.NewMemMapFs(),
			Git: i.git,
			Own: i.own,
			Rep: Repository,
			Ref: ref,
		})
	}

	fnc := func(pat string, fil fs.FileInfo, err error) error {
		{
			if err != nil {
				return tracer.Mask(err)
			}
			if fil.IsDir() {
				return nil
			}
		}

		var ext string
		{
			ext = filepath.Ext(fil.Name())
			if ext != ".yaml" {
				return nil
			}
		}

		var byt []byte
		{
			byt, err = afero.ReadFile(gfs, pat)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		var key string
		{
			key, err = i.envKey(pat)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		var inp *s3.PutObjectInput
		{
			inp = &s3.PutObjectInput{
				Bucket:      aws.String(Bucket),
				Key:         aws.String(key),
				Body:        bytes.NewReader(byt),
				ContentType: aws.String("application/x-yaml"),
			}
		}

		{
			_, err = i.as3.PutObject(context.Background(), inp)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		return nil
	}

	{
		err := afero.Walk(gfs, Directory, fnc)
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
