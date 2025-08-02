package container

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/xh3b4sd/tracer"
)

type detail struct {
	// arn is the well defined Amazon Resource Name of the ECS service.
	arn string
	// clu is the short cluster name that the given service is part of.
	clu string
}

// detail finds all ECS service ARNs that are tagged with the "environment" that
// matches Kayron's runtime configuration. In other words, if Kayron is running
// in "staging", then detail() will find all ECS services labelled with the
// resource tags environment=staging.
func (c *Container) detail() ([]detail, error) {
	var err error

	var inp *resourcegroupstaggingapi.GetResourcesInput
	{
		inp = &resourcegroupstaggingapi.GetResourcesInput{
			ResourceTypeFilters: []string{"ecs:service"},
			TagFilters: []types.TagFilter{
				{
					Key:    aws.String("environment"),
					Values: []string{c.env.Environment},
				},
			},
		}
	}

	var out *resourcegroupstaggingapi.GetResourcesOutput
	{
		out, err = c.tag.GetResources(context.Background(), inp)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var det []detail
	for _, x := range out.ResourceTagMappingList {
		spl := strings.Split(*x.ResourceARN, "/")
		if len(spl) != 3 {
			return nil, tracer.Mask(invalidAmazonResourceNameError, tracer.Context{Key: "arn", Value: *x.ResourceARN})
		}

		det = append(det, detail{
			arn: *x.ResourceARN,
			clu: spl[1],
		})
	}

	return det, nil
}
