package generic

import (
	"fmt"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Generic_Duplicate_string(t *testing.T) {
	testCases := []struct {
		lis []string
		dup []string
	}{
		// Case 000
		{
			lis: []string{},
			dup: nil,
		},
		// Case 001
		{
			lis: []string{
				"55",
				"44",
			},
			dup: nil,
		},
		// Case 002
		{
			lis: []string{
				"33",
				"44",
				"33",
				"33",
			},
			dup: []string{
				"33",
			},
		},
		// Case 003
		{
			lis: []string{
				"33",
				"44",
				"88",
				"22",
				"33",
				"55",
				"66",
				"55",
				"88",
			},
			dup: []string{
				"33",
				"55",
				"88",
			},
		},
		// Case 004
		{
			lis: []string{
				"33",
				"44",
				"88",
				"22",
				"55",
				"66",
			},
			dup: nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			dup := Duplicate(tc.lis)

			slices.Sort(dup)
			slices.Sort(tc.dup)

			if dif := cmp.Diff(tc.dup, dup); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}

func Test_Generic_Duplicate_int64(t *testing.T) {
	testCases := []struct {
		lis []int64
		dup []int64
	}{
		// Case 000
		{
			lis: []int64{},
			dup: nil,
		},
		// Case 001
		{
			lis: []int64{
				55,
				44,
			},
			dup: nil,
		},
		// Case 002
		{
			lis: []int64{
				33,
				44,
				33,
				33,
			},
			dup: []int64{
				33,
			},
		},
		// Case 003
		{
			lis: []int64{
				33,
				44,
				88,
				22,
				33,
				55,
				66,
				55,
				88,
			},
			dup: []int64{
				33,
				55,
				88,
			},
		},
		// Case 004
		{
			lis: []int64{
				33,
				44,
				88,
				22,
				55,
				66,
			},
			dup: nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			dup := Duplicate(tc.lis)

			slices.Sort(dup)
			slices.Sort(tc.dup)

			if dif := cmp.Diff(tc.dup, dup); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
