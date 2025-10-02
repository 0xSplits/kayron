package docker

import (
	"fmt"
	"testing"
)

func Test_Release_Schema_Release_Docker_String_Verify(t *testing.T) {
	testCases := []struct {
		rep String
		mat func(error) bool
	}{
		// Case 000
		{
			rep: "",
			mat: isErr,
		},
		// Case 001
		{
			rep: "UPPER/CASE",
			mat: isErr,
		},
		// Case 002
		{
			rep: "/leading/slash",
			mat: isErr,
		},
		// Case 003
		{
			rep: "trailing/slash/",
			mat: isErr,
		},
		// Case 004
		{
			rep: "double//slash",
			mat: isErr,
		},
		// Case 005
		{
			rep: "has:colon",
			mat: isErr,
		},
		// Case 006
		{
			rep: "has@at",
			mat: isErr,
		},
		// Case 007
		{
			rep: "space inname",
			mat: isErr,
		},
		// Case 008
		{
			rep: "dot..dot",
			mat: isErr,
		},
		// Case 009
		{
			rep: "seg-.-bad",
			mat: isErr,
		},
		// Case 010
		{
			rep: "repo",
			mat: isNil,
		},
		// Case 011
		{
			rep: "library/ubuntu",
			mat: isNil,
		},
		// Case 012
		{
			rep: "my-org/my_app",
			mat: isNil,
		},
		// Case 013
		{
			rep: "a/b-c_d.e",
			mat: isNil,
		},
		// Case 014
		{
			rep: "valid.segment-1/another_segment.2",
			mat: isNil,
		},
		// Case 015
		{
			rep: "verylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylongverylong",
			mat: isNil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			err := tc.rep.Verify()
			if !tc.mat(err) {
				t.Fatal("expected", true, "got", err)
			}
		})
	}
}
