package deploy

import (
	"fmt"
	"testing"

	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/preview"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/release"
	"github.com/0xSplits/kayron/pkg/release/schema/release/deploy/suspend"
)

func Test_Schema_Specification_Service_Deploy_name(t *testing.T) {
	testCases := []struct {
		str Interface
		nam string
	}{
		// Case 000
		{
			str: release.String(""),
			nam: "release",
		},
		// Case 001
		{
			str: suspend.Bool(true),
			nam: "suspend",
		},
		// Case 002
		{
			str: preview.Bool(true),
			nam: "preview",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			nam := name(tc.str)
			if nam != tc.nam {
				t.Fatalf("expected %#v got %#v", tc.nam, nam)
			}
		})
	}
}
