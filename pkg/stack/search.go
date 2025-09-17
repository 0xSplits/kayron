package stack

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/xh3b4sd/tracer"
)

func (s *Stack) Search() (types.Stack, error) {
	var err error

	{
		s.mut.Lock()
		defer s.mut.Unlock()
	}

	if s.sta != nil {
		return *s.sta, nil
	}

	var inp *cloudformation.DescribeStacksInput
	{
		inp = &cloudformation.DescribeStacksInput{
			StackName: aws.String(s.env.CloudformationStack),
		}
	}

	var out *cloudformation.DescribeStacksOutput
	{
		out, err = s.cfc.DescribeStacks(context.Background(), inp)
		if err != nil {
			return types.Stack{}, tracer.Mask(err)
		}
	}

	for _, x := range out.Stacks {
		if !hasEnv(x.Tags, s.env.Environment) {
			continue
		}

		// Only root stacks have no parent stack, so if we find the CloudFormation
		// stack without parent ID, then we found the root stack and return it.

		if !hasPar(x) {
			{
				s.sta = &x
			}

			return x, nil
		}
	}

	return types.Stack{}, tracer.Mask(invalidRootStackError, tracer.Context{Key: "environment", Value: s.env.Environment})
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
