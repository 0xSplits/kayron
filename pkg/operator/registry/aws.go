package registry

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/xh3b4sd/tracer"
)

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

	// Use the DescribeImages API to check whether any given image tag exists in
	// ECR. If we get any error about the image repository or the image tag not
	// existing, then we return false. All other errors need to be propagated back
	// to the caller.

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
