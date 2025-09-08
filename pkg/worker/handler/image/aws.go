package image

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/xh3b4sd/tracer"
)

func (h *Handler) delete(rep string, dig []string) error {
	var err error

	var ids []types.ImageIdentifier
	for _, x := range dig {
		ids = append(ids, types.ImageIdentifier{
			ImageDigest: aws.String(x),
		})
	}

	h.log.Log(
		"level", "info",
		"message", "cleaning up images",
		"amount", strconv.Itoa(len(dig)),
	)

	var inp *ecr.BatchDeleteImageInput
	{
		inp = &ecr.BatchDeleteImageInput{
			RepositoryName: aws.String(rep),
			ImageIds:       ids,
		}
	}

	var out *ecr.BatchDeleteImageOutput
	{
		out, err = h.ecr.BatchDeleteImage(context.Background(), inp)
		if err != nil {
			return tracer.Mask(err, tracer.Context{Key: "repository", Value: rep})
		}
	}

	for _, x := range out.Failures {
		h.log.Log(
			"level", "warning",
			"message", "digest cleanup error",
			"reason", *x.FailureReason,
			"code", string(x.FailureCode),
			"digest", *x.ImageId.ImageDigest,
		)
	}

	return nil
}

func (h *Handler) search(rep string) ([]types.ImageDetail, error) {
	var inp *ecr.DescribeImagesInput
	{
		inp = &ecr.DescribeImagesInput{
			RepositoryName: aws.String(rep),
			MaxResults:     aws.Int32(100),
			Filter: &types.DescribeImagesFilter{
				TagStatus: types.TagStatusTagged,
			},
		}
	}

	var pag *ecr.DescribeImagesPaginator
	{
		pag = ecr.NewDescribeImagesPaginator(h.ecr, inp)
	}

	var out []types.ImageDetail
	for pag.HasMorePages() {
		res, err := pag.NextPage(context.Background())
		if err != nil {
			return nil, tracer.Mask(err, tracer.Context{Key: "repository", Value: rep})
		}

		{
			out = append(out, res.ImageDetails...)
		}
	}

	return out, nil
}
