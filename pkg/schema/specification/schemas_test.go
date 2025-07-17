package specification

import (
	"fmt"
	"testing"

	"github.com/0xSplits/kayron/pkg/schema/specification/infrastructure"
	"github.com/0xSplits/kayron/pkg/schema/specification/service"
	"github.com/0xSplits/kayron/pkg/schema/specification/service/deploy"
	"github.com/0xSplits/kayron/pkg/schema/specification/service/deploy/webhook"
)

// Test_Schemas_Verify_false ensures that invalid schemas can be rejected
// properly, according to their underlying validation rules.
func Test_Schemas_Verify_false(t *testing.T) {
	testCases := []struct {
		sch Schemas
		mat func(error) bool
	}{
		// Case 000, one service, more than one strategy
		{
			sch: Schemas{},
			mat: IsSchemaEmpty,
		},
		// Case 001, one service, more than one strategy
		{
			sch: Schemas{
				{
					Service: service.Services{
						{
							Docker: "kayron",
							Github: "kayron",
						},
					},
				},
				{
					Service: service.Services{},
				},
			},
			mat: IsSchemaEmpty,
		},
		// Case 002, one service, no docker repository
		{
			sch: Schemas{
				{
					Service: service.Services{
						{
							Github: "kayron",
						},
					},
				},
			},
			mat: service.IsServiceDockerEmpty,
		},
		// Case 003, one service, no github repository
		{
			sch: Schemas{
				{
					Service: service.Services{
						{
							Docker: "kayron",
						},
					},
				},
			},
			mat: service.IsServiceGithubEmpty,
		},
		// Case 004, one service, more than one strategy
		{
			sch: Schemas{
				{
					Service: service.Services{
						{
							Docker: "kayron",
							Github: "kayron",
							Deploy: deploy.Deploy{
								Branch:  "feature",
								Release: "v1.8.2",
							},
						},
					},
				},
			},
			mat: deploy.IsDeploymentStrategy,
		},
		// Case 005, many services, more than one strategy
		{
			sch: Schemas{
				{
					Service: service.Services{
						{
							Docker: "kayron",
							Github: "kayron",
							Deploy: deploy.Deploy{
								Release: "v1.8.2",
							},
						},
					},
				},
				{
					Service: service.Services{
						{
							Docker: "specta",
							Github: "specta",
							Deploy: deploy.Deploy{
								Suspend: true,
								Webhook: webhook.Webhooks{
									"POST:https://foo.bar",
								},
							},
						},
					},
				},
				{
					Service: service.Services{
						{
							Docker: "server",
							Github: "server",
							Deploy: deploy.Deploy{
								Branch: "feature",
							},
						},
					},
				},
			},
			mat: deploy.IsDeploymentStrategy,
		},
		// Case 006, duplicated infrastructure shorthands
		{
			sch: Schemas{
				{
					Infrastructure: infrastructure.Infrastructure{
						Shorthand: "prod",
					},
					Service: service.Services{
						{
							Docker: "kayron",
							Github: "kayron",
							Deploy: deploy.Deploy{
								Release: "v1.8.2",
							},
						},
					},
				},
				{
					Infrastructure: infrastructure.Infrastructure{
						Shorthand: "staging",
					},
					Service: service.Services{
						{
							Docker: "specta",
							Github: "specta",
							Deploy: deploy.Deploy{
								Suspend: true,
							},
						},
					},
				},
				{
					Infrastructure: infrastructure.Infrastructure{
						Shorthand: "prod",
					},
					Service: service.Services{
						{
							Docker: "server",
							Github: "server",
							Deploy: deploy.Deploy{
								Branch: "feature",
							},
						},
					},
				},
			},
			mat: IsInfrastructureShorthand,
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

// Test_Schemas_Verify_true ensures that valid schemas can be accepted properly,
// according to their underlying validation rules.
func Test_Schemas_Verify_true(t *testing.T) {
	testCases := []struct {
		sch Schemas
		mat func(error) bool
	}{
		// Case 000, one service, branch strategy
		{
			sch: Schemas{
				{
					Service: service.Services{
						{
							Docker: "kayron",
							Github: "kayron",
							Deploy: deploy.Deploy{
								Branch: "feature",
							},
						},
					},
				},
			},
			mat: deploy.IsDeploymentStrategy,
		},
		// Case 001, one service, release strategy
		{
			sch: Schemas{
				{
					Service: service.Services{
						{
							Docker: "specta",
							Github: "specta",
							Deploy: deploy.Deploy{
								Release: "v1.8.2",
							},
						},
					},
				},
			},
			mat: deploy.IsDeploymentStrategy,
		},
		// Case 002, many services, default and webhook strategy
		{
			sch: Schemas{
				{
					Service: service.Services{
						{
							Docker: "kayron",
							Github: "kayron",
						},
					},
				},
				{
					Service: service.Services{
						{
							Docker: "splits",
							Github: "server",
							Deploy: deploy.Deploy{
								Webhook: webhook.Webhooks{
									"POST:https://foo.bar",
								},
							},
						},
					},
				},
				{
					Service: service.Services{
						{
							Docker: "specta",
							Github: "specta",
						},
					},
				},
			},
			mat: deploy.IsDeploymentStrategy,
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
