package github

import (
	"fmt"
	"strings"
	"testing"
)

func Test_Release_Schema_Release_Github_String_Verify(t *testing.T) {
	testCases := []struct {
		rep String
		mat func(error) bool
	}{
		{
			// Case 000, simple
			rep: "repo",
			mat: isNil,
		},
		{
			// Case 001, mix of cases, digits, separators
			rep: "Repo_123.name",
			mat: isNil,
		},
		{
			// Case 002, max length (100)
			rep: String(strings.Repeat("a", 100)),
			mat: isNil,
		},

		{
			// Case 003, empty
			rep: "",
			mat: isErr,
		},
		{
			// Case 004, starts with dot
			rep: ".repo",
			mat: isErr,
		},
		{
			// Case 005, ends with .git
			rep: "foo.git",
			mat: isErr,
		},
		{
			// Case 006, over max length
			rep: String(strings.Repeat("a", 101)),
			mat: isErr,
		},
		{
			// Case 007, slash
			rep: "re/po",
			mat: isErr,
		},
		{
			// Case 008, punctuation not allowed
			rep: "foo!bar",
			mat: isErr,
		},
		{
			// Case 009, non-ascii
			rep: "schröder",
			mat: isErr,
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
