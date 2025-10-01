package schema

import (
	"fmt"
	"testing"

	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/preview"
	"github.com/0xSplits/kayron/pkg/release/schema/release/docker"
	"github.com/0xSplits/kayron/pkg/release/schema/release/github"
	"github.com/0xSplits/kayron/pkg/release/schema/release/labels"
	"github.com/0xSplits/kayron/pkg/release/schema/release/provider"
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
			mat: release.IsServiceDefinitionEmpty,
		},
		// Case 001, no service
		{
			sch: Schema{
				Release: release.Slice{},
			},
			mat: release.IsServiceDefinitionEmpty,
		},
		// Case 002, one service, no deployment strategy
		{
			sch: Schema{
				Release: release.Slice{
					{
						Docker: docker.String("kayron"),
						Github: github.String("kayron"),
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
			mat: release.IsServiceDeployEmpty,
		},
		// Case 003, one service, no deployment strategy, cloudformation
		{
			sch: Schema{
				Release: release.Slice{
					{
						Github:   github.String("infrastructure"),
						Provider: provider.String("cloudformation"),
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
			mat: release.IsServiceDeployEmpty,
		},
		// Case 004, one service, no docker repository, no provider setting
		{
			sch: Schema{
				Release: release.Slice{
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
			mat: release.IsServiceProviderEmpty,
		},
		// Case 005, one service, no github repository
		{
			sch: Schema{
				Release: release.Slice{
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
			mat: release.IsServiceGithubEmpty,
		},
		// Case 006, one service, more than one strategy
		{
			sch: Schema{
				Release: release.Slice{
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
				Release: release.Slice{
					{
						Docker: docker.String("kayron"),
						Github: github.String("kayron"),
						Deploy: deploy.Struct{
							Release: "v1.8.2",
							Suspend: true,
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
				Release: release.Slice{
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
			mat: release.IsServiceLabelsEmpty,
		},
		// Case 009, one provider, preview deployments
		{
			sch: Schema{
				Release: release.Slice{
					{
						Github:   github.String("infrastructure"),
						Provider: provider.String("cloudformation"),
						Deploy: deploy.Struct{
							Preview: preview.Bool(true),
						},
						Labels: labels.Struct{
							Source: "foo",
						},
					},
				},
			},
			mat: release.IsServiceDeployPreview,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			err := tc.sch.Verify()
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
				Release: release.Slice{
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
				Release: release.Slice{
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
