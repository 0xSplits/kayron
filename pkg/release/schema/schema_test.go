package schema

import (
	"fmt"
	"testing"

	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/0xSplits/kayron/pkg/release/schema/service/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/service/deploy/webhook"
	"github.com/0xSplits/kayron/pkg/release/schema/service/docker"
	"github.com/0xSplits/kayron/pkg/release/schema/service/github"
	"github.com/0xSplits/kayron/pkg/release/schema/service/labels"
	"github.com/0xSplits/kayron/pkg/release/schema/service/provider"
)

// Test_Schema_Verify_failure ensures that invalid service definitions are
// rejected properly, according to their underlying validation rules.
func Test_Schema_Verify_failure(t *testing.T) {
	testCases := []struct {
		sch Schema
		mat func(error) bool
	}{
		// Case 000, no service, nil
		{
			sch: Schema{},
			mat: service.IsServiceDefinitionEmpty,
		},
		// Case 001, no service
		{
			sch: Schema{
				Service: service.Slice{},
			},
			mat: service.IsServiceDefinitionEmpty,
		},
		// Case 002, one service, no deployment strategy
		{
			sch: Schema{
				Service: service.Slice{
					{
						Docker: docker.String("kayron"),
						Github: github.String("kayron"),
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
			mat: service.IsServiceDeployEmpty,
		},
		// Case 003, one service, no deployment strategy, cloudformation
		{
			sch: Schema{
				Service: service.Slice{
					{
						Github:   github.String("infrastructure"),
						Provider: provider.String("cloudformation"),
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
			mat: service.IsServiceDeployEmpty,
		},
		// Case 004, one service, no docker repository, no provider setting
		{
			sch: Schema{
				Service: service.Slice{
					{
						Github: "kayron",
						Deploy: deploy.Struct{
							Release: "v1.8.2",
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
			mat: service.IsServiceProviderEmpty,
		},
		// Case 005, one service, no github repository
		{
			sch: Schema{
				Service: service.Slice{
					{
						Docker: "kayron",
						Deploy: deploy.Struct{
							Release: "v1.8.2",
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
			mat: service.IsServiceGithubEmpty,
		},
		// Case 006, one service, more than one strategy
		{
			sch: Schema{
				Service: service.Slice{
					{
						Docker: docker.String("kayron"),
						Github: github.String("kayron"),
						Deploy: deploy.Struct{
							Branch:  "feature",
							Release: "v1.8.2",
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
			mat: deploy.IsDeploymentStrategy,
		},
		// Case 007, many services, more than one strategy
		{
			sch: Schema{
				Service: service.Slice{
					{
						Docker: docker.String("kayron"),
						Github: github.String("kayron"),
						Deploy: deploy.Struct{
							Release: "v1.8.2",
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
					{
						Docker: docker.String("specta"),
						Github: github.String("specta"),
						Deploy: deploy.Struct{
							Suspend: true,
							Webhook: webhook.Slice{
								"POST:https://foo.bar",
							},
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
					{
						Docker: docker.String("server"),
						Github: github.String("server"),
						Deploy: deploy.Struct{
							Branch: "feature",
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
			mat: deploy.IsDeploymentStrategy,
		},
		// Case 008, many services, no labels
		{
			sch: Schema{
				Service: service.Slice{
					{
						Docker: docker.String("kayron"),
						Github: github.String("kayron"),
						Deploy: deploy.Struct{
							Release: "v1.8.2",
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
					{
						Docker: docker.String("specta"),
						Github: github.String("specta"),
						Deploy: deploy.Struct{
							Suspend: true,
						},
					},
					{
						Docker: docker.String("server"),
						Github: github.String("server"),
						Deploy: deploy.Struct{
							Branch: "feature",
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
			mat: service.IsServiceLabelsEmpty,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			err := tc.sch.Verify()
			fmt.Printf("%#v\n", err)
			if !tc.mat(err) {
				t.Fatal("expected", true, "got", err)
			}
		})
	}
}

// Test_Schema_Verify_success ensures that invalid service definitions are
// accepted properly, according to their underlying validation rules.
func Test_Schema_Verify_success(t *testing.T) {
	testCases := []struct {
		sch Schema
	}{
		// Case 000, one service, branch strategy
		{
			sch: Schema{
				Service: service.Slice{
					{
						Docker: docker.String("kayron"),
						Github: github.String("kayron"),
						Deploy: deploy.Struct{
							Branch: "feature",
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
		},
		// Case 001, one service, release strategy
		{
			sch: Schema{
				Service: service.Slice{
					{
						Docker: docker.String("specta"),
						Github: github.String("specta"),
						Deploy: deploy.Struct{
							Release: "v1.8.2",
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
		},
		// Case 002, many services, webhook strategy
		{
			sch: Schema{
				Service: service.Slice{
					{
						Docker: "kayron",
						Github: "kayron",
						Deploy: deploy.Struct{
							Release: "v1.8.2",
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
					{
						Docker: docker.String("splits"),
						Github: github.String("server"),
						Deploy: deploy.Struct{
							Webhook: webhook.Slice{
								"POST:https://foo.bar",
							},
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
					{
						Docker: docker.String("specta"),
						Github: github.String("specta"),
						Deploy: deploy.Struct{
							Release: "v1.8.2",
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			err := tc.sch.Verify()
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
		})
	}
}
