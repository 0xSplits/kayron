package cache

import (
	"fmt"
	"testing"

	"github.com/0xSplits/kayron/pkg/release/schema/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/docker"
	"github.com/0xSplits/kayron/pkg/release/schema/release/github"
)

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
