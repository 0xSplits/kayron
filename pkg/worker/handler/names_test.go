package handler

import (
	"fmt"
	"slices"
	"testing"

	"github.com/0xSplits/kayron/pkg/worker/handler/operator"
)

func Test_Worker_Handler_Names(t *testing.T) {
	testCases := []struct {
		han []Interface
		nam []string
	}{
		// Case 000
		{
			han: []Interface{
				&operator.Handler{},
			},
			nam: []string{
				"operator",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			nam := Names(tc.han)
			if !slices.Equal(nam, tc.nam) {
				t.Fatalf("expected %#v got %#v", tc.nam, nam)
			}
		})
	}
}
