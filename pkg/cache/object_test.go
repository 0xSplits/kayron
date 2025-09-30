package cache

import (
	"fmt"
	"testing"

	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/artifact/condition"
	"github.com/0xSplits/kayron/pkg/release/artifact/reference"
	"github.com/0xSplits/kayron/pkg/release/artifact/scheduler"
	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/suspend"
	"github.com/0xSplits/kayron/pkg/release/schema/release/docker"
	"github.com/0xSplits/kayron/pkg/release/schema/release/github"
	"github.com/google/go-cmp/cmp"
)

func Test_Cache_Object_Drift_ready(t *testing.T) {
	testCases := []struct {
		rel Object
		dft bool
	}{
		// Case 000, empty release, no drift
		{
			rel: Object{},
			dft: false,
		},
		// Case 001, no drift
		{
			rel: Object{
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
			rel: Object{
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
			rel: Object{
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
			rel: Object{
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
			dft := tc.rel.Drift(true)
			if dif := cmp.Diff(tc.dft, dft); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}

func Test_Cache_Object_Drift_waiting(t *testing.T) {
	testCases := []struct {
		rel Object
		dft bool
	}{
		// Case 000, empty release, no drift
		{
			rel: Object{},
			dft: false,
		},
		// Case 001, no drift
		{
			rel: Object{
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
			rel: Object{
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
			dft: true,
		},
		// Case 003, drift
		{
			rel: Object{
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
			rel: Object{
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
			dft := tc.rel.Drift(false)
			if dif := cmp.Diff(tc.dft, dft); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}

func Test_Cache_Object_Parameter(t *testing.T) {
	testCases := []struct {
		obj Object
		par string
	}{
		// Case 000
		{
			obj: Object{
				Release: release.Struct{Github: github.String("infrastructure")},
				kin:     Infrastructure,
			},
			par: "InfrastructureVersion",
		},
		// Case 001
		{
			obj: Object{
				Release: release.Struct{Github: github.String("heLLowOrlD")},
				kin:     Infrastructure,
			},
			par: "HelloworldVersion",
		},
		// Case 002
		{
			obj: Object{
				Release: release.Struct{Docker: docker.String("specta")},
				kin:     Service,
			},
			par: "SpectaVersion",
		},
		// Case 003
		{
			obj: Object{
				Release: release.Struct{Docker: docker.String("fOobAR")},
				kin:     Service,
			},
			par: "FoobarVersion",
		},
		// Case 004, with dash
		{
			obj: Object{
				Release: release.Struct{Docker: docker.String("splits-lite")},
				kin:     Service,
			},
			par: "SplitsLiteVersion",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			par := tc.obj.Parameter()
			if par != tc.par {
				t.Fatalf("expected %#v got %#v", tc.par, par)
			}
		})
	}
}
