package roghfs

import (
	"fmt"
	"testing"
)

func Test_Roghfs_Parse_failure(t *testing.T) {
	testCases := []struct {
		str string
		mat func(error) bool
	}{
		// Case 000
		{
			str: "",
			mat: IsInvalidRepositoryUrl,
		},
		// Case 001
		{
			str: "   ",
			mat: IsInvalidRepositoryUrl,
		},
		// Case 002
		{
			str: "http://github.com/owner/name",
			mat: IsInvalidRepositoryUrl,
		},
		// Case 003
		{
			str: "https://example.com/foo/bar",
			mat: IsInvalidRepositoryUrl,
		},
		// Case 004
		{
			str: "https://github.com/owner",
			mat: IsInvalidRepositoryUrl,
		},
		// Case 005
		{
			str: "https://github.com/owner/name/extra",
			mat: IsInvalidRepositoryUrl,
		},
		// Case 006
		{
			str: "https://github.com//name",
			mat: IsInvalidRepositoryUrl,
		},
		// Case 007
		{
			str: "https://github.com/owner/",
			mat: IsInvalidRepositoryUrl,
		},
		// Case 008
		{
			str: "https://github.com/owner/name?foo=bar",
			mat: IsInvalidRepositoryUrl,
		},
		// Case 009
		{
			str: "https://github.com/owner/name#frag",
			mat: IsInvalidRepositoryUrl,
		},
		// Case 010
		{
			str: "https://github.com/alpha/one.git",
			mat: IsInvalidRepositoryUrl,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			_, _, err := Parse(tc.str)
			if !tc.mat(err) {
				t.Fatal("expected", true, "got", err)
			}
		})
	}
}

func Test_Roghfs_Parse_success(t *testing.T) {
	testCases := []struct {
		str string
		own string
		rep string
	}{
		// Case 000
		{
			str: "https://github.com/alpha/one",
			own: "alpha",
			rep: "one",
		},
		// Case 001
		{
			str: "https://github.com/beta/two/",
			own: "beta",
			rep: "two",
		},
		// Case 002
		{
			str: "  https://github.com/gamma/three  ",
			own: "gamma",
			rep: "three",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			own, rep, err := Parse(tc.str)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
			if own != tc.own {
				t.Fatal("expected", tc.own, "got", own)
			}
			if rep != tc.rep {
				t.Fatal("expected", tc.rep, "got", rep)
			}
		})
	}
}
