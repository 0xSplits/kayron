package specification

import (
	"fmt"
	"testing"

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
							GitHub: "kayron",
						},
					},
				},
				{
					Service: service.Services{},
				},
			},
			mat: IsSchemaEmpty,
		},
		// Case 002, one service, more than one strategy
		{
			sch: Schemas{
				{
					Service: service.Services{
						{
							Docker: "kayron",
							GitHub: "kayron",
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
		// Case 003, many services, more than one strategy
		{
			sch: Schemas{
				{
					Service: service.Services{
						{
							Docker: "kayron",
							GitHub: "kayron",
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
							GitHub: "specta",
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
							GitHub: "server",
							Deploy: deploy.Deploy{
								Branch: "feature",
							},
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
							GitHub: "kayron",
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
							GitHub: "specta",
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
							GitHub: "kayron",
						},
					},
				},
				{
					Service: service.Services{
						{
							Docker: "splits",
							GitHub: "server",
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
							GitHub: "specta",
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
