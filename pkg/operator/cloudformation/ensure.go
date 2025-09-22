package cloudformation

import (
	"context"
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/xh3b4sd/tracer"
)

func (c *CloudFormation) Ensure() error {
	var err error

	var tag []types.Tag
	for k, v := range c.env.CloudformationTags {
		tag = append(tag, types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	var url string
	{
		url = c.temUrl()
	}

	c.log.Log(
		"level", "info",
		"message", "updating cloudformation stack",
		"name", c.env.CloudformationStack,
		"url", url,
	)

	// Make sure we respect the dry run flag when attempting to update any stack
	// within CloudFormation, because "dry run" effectively means "read only". So
	// if the dry run flag is set in e.g. the operator's integration test, then we
	// want to emit the logs, but prevent making the network calls.

	var inp *cloudformation.UpdateStackInput
	{
		inp = &cloudformation.UpdateStackInput{
			StackName:   aws.String(c.env.CloudformationStack),
			TemplateURL: aws.String(url),
			Parameters:  c.temPar(c.cac.Releases()),
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
	//     emit deployment event via Specta, if updated
	//

	return nil
}

func (c *CloudFormation) temPar(rel []cache.Object) []types.Parameter {
	var par []types.Parameter

	for k, v := range c.env.CloudformationParameters {
		par = append(par, types.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}

	// Inject all desired artifact versions into the parameters that we are just
	// about to deploy, but only for main release definitions, not for auto
	// generated preview deployments. Injecting the template parameters after all
	// user inputs have been applied above guarantees that only the release
	// versions as defined in the release source repository will ever be applied.

	for _, x := range rel {
		if x.Preview() {
			continue
		}

		par = append(par, types.Parameter{
			ParameterKey:   aws.String(x.Parameter()),
			ParameterValue: aws.String(x.Version()),
		})
	}

	return par
}

// temUrl returns the environment specific template URL for the root stack
// CloudFormation template.
func (c *CloudFormation) temUrl() string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s/index.yaml", c.env.S3Bucket, c.cfc.Options().Region, c.env.Environment)
}
