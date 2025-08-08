package template

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/xh3b4sd/tracer"
)

// rooSta finds the first root level CloudFormation stack that is tagged with
// the "environment" that matches Kayron's runtime configuration. In other
// words, if Kayron is running in "staging", then rooSta will find the first
// CloudFormation root stack labelled with the resource tags
// environment=staging.
func (t *Template) rooSta() (types.Stack, error) {
	var err error

	// Create a new paginator for each reconciliation loop in order to guarantee
	// data integrity per iteration.
	var pag *cloudformation.DescribeStacksPaginator
	{
		pag = cloudformation.NewDescribeStacksPaginator(t.cfc, &cloudformation.DescribeStacksInput{})
	}

	for pag.HasMorePages() {
		var out *cloudformation.DescribeStacksOutput
		{
			out, err = pag.NextPage(context.Background())
			if err != nil {
				return types.Stack{}, tracer.Mask(err)
			}
		}

		for _, x := range out.Stacks {
			if !hasEnv(x.Tags, t.env.Environment) {
				continue
			}

			// Only root stacks have no parent stack, so if we find the CloudFormation
			// stack without parent ID, then we found the root stack and return it.

			if !hasPar(x) {
				return x, nil
			}
		}
	}

	return types.Stack{}, tracer.Mask(missingRootStackError, tracer.Context{Key: "environment", Value: t.env.Environment})
}

func hasEnv(tags []types.Tag, env string) bool {
	for _, t := range tags {
		if t.Key != nil && t.Value != nil && *t.Key == "environment" && *t.Value == env {
			return true
		}
	}

	return false
}

func hasPar(sta types.Stack) bool {
	return sta.ParentId == nil
}
