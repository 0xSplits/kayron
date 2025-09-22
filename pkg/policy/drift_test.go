package policy

import (
	"fmt"
	"testing"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/artifact/condition"
	"github.com/0xSplits/kayron/pkg/release/artifact/reference"
	"github.com/0xSplits/kayron/pkg/release/artifact/scheduler"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/suspend"
	"github.com/google/go-cmp/cmp"
)

func Test_Operator_Policy_drift(t *testing.T) {
	testCases := []struct {
		rel cache.Object
		dft bool
	}{
		// Case 000, empty release, no drift
		{
			rel: cache.Object{},
			dft: false,
		},
		// Case 001, no drift
		{
			rel: cache.Object{
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
			dft: false,
		},
		// Case 002, drift, failed condition
		{
			rel: cache.Object{
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
			dft: false,
		},
		// Case 003, drift
		{
			rel: cache.Object{
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
			dft: true,
		},
		// Case 004, drift, suspended
		{
			rel: cache.Object{
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
			dft: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			dft := drift(tc.rel)
			if dif := cmp.Diff(tc.dft, dft); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
