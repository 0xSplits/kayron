package generic

import (
	"fmt"
	"testing"
)

func Test_Generic_Duplicate_string(t *testing.T) {
	testCases := []struct {
		lis []string
		dup bool
	}{
		// Case 000
		{
			lis: []string{},
			dup: false,
		},
		// Case 001
		{
			lis: []string{
				"55",
				"44",
			},
			dup: false,
		},
		// Case 002
		{
			lis: []string{
				"33",
				"44",
				"33",
				"33",
			},
			dup: true,
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
			dup: true,
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
			dup: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			dup := Duplicate(tc.lis)
			if dup != tc.dup {
				t.Fatalf("expected %#v got %#v", tc.dup, dup)
			}
		})
	}
}

func Test_Generic_Duplicate_int64(t *testing.T) {
	testCases := []struct {
		lis []int64
		dup bool
	}{
		// Case 000
		{
			lis: []int64{},
			dup: false,
		},
		// Case 001
		{
			lis: []int64{
				55,
				44,
			},
			dup: false,
		},
		// Case 002
		{
			lis: []int64{
				33,
				44,
				33,
				33,
			},
			dup: true,
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
			dup: true,
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
			dup: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			dup := Duplicate(tc.lis)
			if dup != tc.dup {
				t.Fatalf("expected %#v got %#v", tc.dup, dup)
			}
		})
	}
}
