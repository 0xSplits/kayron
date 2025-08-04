package registry

import (
	"context"
	"strconv"

	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

func (r *Registry) Ensure() error {
	var err error

	var ser []service.Service
	for i := range r.ser.Length() {
		var s service.Service
		{
			s, _ = r.ser.Search(i)
		}

		{
			ser = append(ser, s)
		}
	}

	// Check whether the desired Docker image exists within the underlying
	// container registry, if the current and desired state differs.

	fnc := func(i int, x service.Service) error {
		var err error

		var cur string
		var des string
		{
			cur, _ = r.art.Search(artifact.ReferenceCurrent(i))
			des, _ = r.art.Search(artifact.ReferenceDesired(i))
		}

		// We do not have to do any work here if the currently deployed service
		// already matches the desired service release. Note that we also have to
		// prevent the infrastructure configuration that is part of the cached
		// services to initiate its own image lookup, because there is no image for
		// our templates. The current setup is not optimal and a better data
		// structure may be used in a future refactoring, so that we do not have to
		// do these provider checks anymore.

		if cur == des || x.Provider == "cloudformation" {
			return nil
		}

		var exi bool
		{
			exi, err = r.imaExi(x.Docker.String(), des)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		var str string
		{
			str = strconv.FormatBool(exi)
		}

		{
			r.art.Update(artifact.ContainerExists(i), str)
		}

		r.log.Log(
			"level", "debug",
			"message", "executed image check",
			"image", x.Docker.String(),
			"tag", des,
			"exists", str,
		)

		return nil
	}

	{
		err = parallel.Slice(ser, fnc)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (r *Registry) imaExi(rep string, tag string) (bool, error) {
	var err error

	var inp *ecr.DescribeImagesInput
	{
		inp = &ecr.DescribeImagesInput{
			RepositoryName: aws.String(rep),
			ImageIds: []types.ImageIdentifier{
				{ImageTag: aws.String(tag)},
			},
		}
	}

	{
		_, err = r.ecr.DescribeImages(context.Background(), inp)
		if isImageNotFound(err) {
			return false, nil
		} else if err != nil {
			return false, tracer.Mask(err)
		}
	}

	return true, nil
}
