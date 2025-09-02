package cloudformation

import (
	"fmt"
	"testing"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/envvar"
	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/artifact/condition"
	"github.com/0xSplits/kayron/pkg/release/artifact/reference"
	"github.com/0xSplits/kayron/pkg/release/artifact/scheduler"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/suspend"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Test_Operator_CloudFormation_temPar ensures that we do not use new release
// versions for release artifacts that are e.g. missing a pushed docker image,
// or have their deployment strategy suspended.
func Test_Operator_CloudFormation_temPar(t *testing.T) {
	testCases := []struct {
		rel []cache.Object
		par []types.Parameter
	}{
		// Case 000, no release
		{
			rel: []cache.Object{},
			par: nil,
		},
		// Case 001, existing release
		{
			rel: []cache.Object{
				{
					Artifact: artifact.Struct{
						Condition: condition.Struct{
							Success: true,
						},
						Scheduler: scheduler.Struct{
							Current: "foo",
						},
						Reference: reference.Struct{
							Desired: "foo",
						},
					},
				},
			},
			par: []types.Parameter{
				{ParameterKey: aws.String("Version"), ParameterValue: aws.String("foo")},
			},
		},
		// Case 002, new release not ready
		{
			rel: []cache.Object{
				{
					Artifact: artifact.Struct{
						Condition: condition.Struct{
							Success: false,
						},
						Scheduler: scheduler.Struct{
							Current: "foo",
						},
						Reference: reference.Struct{
							Desired: "bar",
						},
					},
				},
			},
			par: []types.Parameter{
				{ParameterKey: aws.String("Version"), ParameterValue: aws.String("foo")},
			},
		},
		// Case 003, new release ready
		{
			rel: []cache.Object{
				{
					Artifact: artifact.Struct{
						Condition: condition.Struct{
							Success: true,
						},
						Scheduler: scheduler.Struct{
							Current: "foo",
						},
						Reference: reference.Struct{
							Desired: "bar",
						},
					},
				},
			},
			par: []types.Parameter{
				{ParameterKey: aws.String("Version"), ParameterValue: aws.String("bar")},
			},
		},
		// Case 004, new release ready but suspended
		{
			rel: []cache.Object{
				{
					Artifact: artifact.Struct{
						Condition: condition.Struct{
							Success: true,
						},
						Scheduler: scheduler.Struct{
							Current: "foo",
						},
						Reference: reference.Struct{
							Desired: "bar",
						},
					},
					Release: release.Struct{
						Deploy: deploy.Struct{
							Suspend: suspend.Bool(true),
						},
					},
				},
			},
			par: []types.Parameter{
				{ParameterKey: aws.String("Version"), ParameterValue: aws.String("foo")},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var clo *CloudFormation
			{
				clo = &CloudFormation{}
			}

			var opt []cmp.Option
			{
				opt = []cmp.Option{
					cmpopts.IgnoreUnexported(types.Parameter{}),
				}
			}

			par := clo.temPar(tc.rel)
			if dif := cmp.Diff(tc.par, par, opt...); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}

// Test_Operator_CloudFormation_temUrl verifies that the actually injected
// environment is properly used for the template URL of the root stack
// CloudFormation template.
func Test_Operator_CloudFormation_temUrl(t *testing.T) {
	testCases := []struct {
		buc string
		reg string
		env string
		url string
	}{
		// Case 000
		{
			buc: "splits-cf-templates",
			reg: "us-west-2",
			env: "testing",
			url: "https://splits-cf-templates.s3.us-west-2.amazonaws.com/testing/index.yaml",
		},
		// Case 001
		{
			buc: "template-files",
			reg: "us-east-1",
			env: "staging",
			url: "https://template-files.s3.us-east-1.amazonaws.com/staging/index.yaml",
		},
		// Case 002
		{
			buc: "splits",
			reg: "eu-central-1",
			env: "production",
			url: "https://splits.s3.eu-central-1.amazonaws.com/production/index.yaml",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var clo *CloudFormation
			{
				clo = &CloudFormation{
					cfc: cloudformation.NewFromConfig(aws.Config{
						Region: tc.reg,
					}),
					env: envvar.Env{
						Environment: tc.env,
						S3Bucket:    tc.buc,
					},
				}
			}

			url := clo.temUrl()
			if dif := cmp.Diff(tc.url, url); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
