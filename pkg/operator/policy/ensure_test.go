package policy

import (
	"fmt"
	"testing"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/cancel"
	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/artifact/condition"
	"github.com/0xSplits/kayron/pkg/release/artifact/reference"
	"github.com/0xSplits/kayron/pkg/release/artifact/scheduler"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/suspend"
	"github.com/xh3b4sd/logger"
)

func Test_Operator_Policy_Ensure(t *testing.T) {
	testCases := []struct {
		rel []cache.Object
		mat func(error) bool
	}{
		// Case 000, no release should cancel
		{
			rel: []cache.Object{},
			mat: cancel.Is,
		},
		// Case 001, single release with no state drift should cancel
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
			mat: cancel.Is,
		},
		// Case 002, single release with state drift and failed condition should cancel
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
			mat: cancel.Is,
		},
		// Case 003, single release with state drift should not cancel
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
			mat: isNil,
		},
		// Case 004, many releases, one with state drift, should not cancel
		{
			rel: []cache.Object{
				// cancel
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

				// cancel
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

				// update
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
			mat: isNil,
		},
		// Case 005, single release with state drift but suspended should cancel
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
			mat: cancel.Is,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var cac *cache.Cache
			{
				cac = cache.New(cache.Config{
					Log: logger.Fake(),
				})
			}

			var log logger.Interface
			{
				log = logger.Fake()
			}

			var pol *Policy
			{
				pol = New(Config{
					Cac: cac,
					Log: log,
				})
			}

			err := pol.ensure(tc.rel)
			if !tc.mat(err) {
				t.Fatal("expected", true, "got", err)
			}
		})
	}
}
