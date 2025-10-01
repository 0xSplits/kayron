package release

import (
	"fmt"
	"testing"
)

func Test_Release_Schema_Release_Deploy_Release_Verify(t *testing.T) {
	testCases := []struct {
		tag String
		mat func(error) bool
	}{
		// Case 000
		{
			tag: "",
			mat: isErr,
		},
		// Case 001
		{
			tag: "1.2.3",
			mat: isErr,
		},
		// Case 002
		{
			tag: "v1.2.x",
			mat: isErr,
		},
		// Case 003
		{
			tag: "v1.2.3-",
			mat: isErr,
		},
		// Case 004
		{
			tag: "v1.2.3+build",
			mat: isErr,
		},
		// Case 005
		{
			tag: "v1.2.3-abc+def",
			mat: isErr,
		},
		// Case 006
		{
			tag: "v1.2.3-alpha beta",
			mat: isErr,
		},
		// Case 007
		{
			tag: "v1.2.3-a/b",
			mat: isErr,
		},
		// Case 008
		{
			tag: "v1.2.3-@abc",
			mat: isErr,
		},
		// Case 009
		{
			tag: "v10.20.30-ABC_ok.123-xyz",
			mat: isErr,
		},

		// Case 010
		{
			tag: "v1.2",
			mat: isNil,
		},
		// Case 011
		{
			tag: "v0.1.0",
			mat: isNil,
		},
		// Case 012
		{
			tag: "v1.8.2",
			mat: isNil,
		},
		// Case 013
		{
			tag: "v1.8.3-ffce1e2",
			mat: isNil,
		},
		// Case 014
		{
			tag: "v1.2.3-abc.123",
			mat: isNil,
		},
		// Case 015
		{
			tag: "v0.0.0",
			mat: isNil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			err := tc.tag.Verify()
			if !tc.mat(err) {
				t.Fatal("expected", true, "got", err)
			}
		})
	}
}
