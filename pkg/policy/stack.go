package policy

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/xh3b4sd/tracer"
)

// Stack returns the state of the configured stack object and caches the first
// valid search result so that consecutive executions of Stack prevent network
// calls. This behaviour guarantees consistent stack object state within
// reconciliation loops.
func (p *Policy) Stack() (types.Stack, error) {
	var err error

	{
		p.mut.Lock()
		defer p.mut.Unlock()
	}

	if p.sta != nil {
		return *p.sta, nil
	}

	var inp *cloudformation.DescribeStacksInput
	{
		inp = &cloudformation.DescribeStacksInput{
			StackName: aws.String(p.env.CloudformationStack),
		}
	}

	var out *cloudformation.DescribeStacksOutput
	{
		out, err = p.cfc.DescribeStacks(context.Background(), inp)
		if err != nil {
			return types.Stack{}, tracer.Mask(err)
		}
	}

	for _, x := range out.Stacks {
		if !hasEnv(x.Tags, p.env.Environment) {
			continue
		}

		// Only root stacks have no parent stack, so if we find the CloudFormation
		// stack without parent ID, then we found the root stack and return it.

		if !hasPar(x) {
			{
				p.sta = &x
			}

			return x, nil
		}
	}

	return types.Stack{}, tracer.Mask(invalidRootStackError, tracer.Context{Key: "environment", Value: p.env.Environment})
}

func hasEnv(tags []types.Tag, env string) bool {
	for _, x := range tags {
		if aws.ToString(x.Key) == "environment" && aws.ToString(x.Value) == env {
			return true
		}
	}

	return false
}

// hasPar returns true if the given stack has a parent ID that is not nil.
func hasPar(sta types.Stack) bool {
	return sta.ParentId != nil
}
