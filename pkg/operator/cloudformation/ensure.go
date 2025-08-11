package cloudformation

import (
	"context"
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/constant"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/xh3b4sd/tracer"
)

func (c *CloudFormation) Ensure() error {
	var err error

	var inf cache.Object
	{
		inf = c.cac.Infrastructure()
	}

	var par []types.Parameter
	for k, v := range c.env.CloudformationParameters {
		par = append(par, types.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}

	// Inject the new template version into the parameters that we are just about
	// to deploy. Injecting this parameter after all user inputs have been applied
	// above guarantees that only the infrastructure version as defined with its
	// own infrastructure release will ever be applied.

	{
		par = append(par, types.Parameter{
			ParameterKey:   aws.String(constant.KayronTemplateVersion),
			ParameterValue: aws.String(inf.Artifact.Reference.Desired),
		})
	}

	var tag []types.Tag
	for k, v := range c.env.CloudformationTags {
		tag = append(tag, types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	c.log.Log(
		"level", "debug",
		"message", "updating cloudformation stack",
		"name", c.env.CloudformationStack,
	)

	// Make sure we respect the dry run flag when attempting to update any stack
	// within CloudFormation, because "dry run" effectively means "read only". So
	// if the dry run flag is set in e.g. the operator's integration test, then we
	// want to emit the logs, but prevent making the network calls.

	var inp *cloudformation.UpdateStackInput
	{
		inp = &cloudformation.UpdateStackInput{
			StackName:   aws.String(c.env.CloudformationStack),
			TemplateURL: aws.String(c.temUrl()),
			Parameters:  par,
			Capabilities: []types.Capability{
				types.CapabilityCapabilityIam,
			},
			Tags: tag,
		}
	}

	if !c.dry {
		_, err = c.cfc.UpdateStack(context.Background(), inp)
		if isNoStateDriftError(err) {
			c.log.Log(
				"level", "debug",
				"message", "no state drift",
			)
		} else if err != nil {
			return tracer.Mask(err)
		}
	}

	// TODO
	//
	//     add inhibition operator function to cancel reconciliation if stack is updating
	//     emit deployment event, if updated
	//

	return nil
}

func (c *CloudFormation) temUrl() string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/testing/index.yaml", c.env.S3Bucket, c.cfc.Options().Region)
}
