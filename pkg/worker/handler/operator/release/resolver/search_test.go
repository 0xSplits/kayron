package resolver

import (
	"fmt"
	"testing"
)

func Test_Worker_Handler_Releases_Resolver_Search_failure(t *testing.T) {
	testCases := []struct {
		env string
		com []func(string) (string, error)
		lat func() (string, error)
		mat func(error) bool
	}{
		// Case 000, production, no release
		{
			env: "production",
			com: nil,
			lat: func() (string, error) { return "", releaseNotFoundError },
			mat: IsReleaseNotFound,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			_, err := Search(&fakeResolver{0, tc.com, tc.lat}, tc.env)
			if !tc.mat(err) {
				t.Fatal("expected", true, "got", err)
			}
		})
	}
}

func Test_Worker_Handler_Releases_Resolver_Search_success(t *testing.T) {
	testCases := []struct {
		env string
		com []func(string) (string, error)
		lat func() (string, error)
		ref string
	}{
		// Case 000, staging, default branch
		{
			env: "staging",
			com: []func(string) (string, error){
				func(string) (string, error) { return "1234", nil },
			},
			lat: nil,
			ref: "1234",
		},
		// Case 001, production, release tag
		{
			env: "production",
			com: nil,
			lat: func() (string, error) { return "v0.1.0", nil },
			ref: "v0.1.0",
		},
		// Case 002, testing, commit sha
		{
			env: "testing",
			com: []func(string) (string, error){
				func(string) (string, error) { return "5678", nil },
			},
			lat: nil,
			ref: "5678",
		},
		// Case 003, testing, no sha, default branch
		{
			env: "testing",
			com: []func(string) (string, error){
				func(string) (string, error) { return "", nil },
				func(string) (string, error) { return "4321", nil },
			},
			lat: nil,
			ref: "4321",
		},
		// Case 004, arbitrary test environment, commit sha
		{
			env: "melissa",
			com: []func(string) (string, error){
				func(string) (string, error) { return "5678", nil },
			},
			lat: nil,
			ref: "5678",
		},
		// Case 005, arbitrary test environment, no sha, default branch
		{
			env: "melissa",
			com: []func(string) (string, error){
				func(string) (string, error) { return "", nil },
				func(string) (string, error) { return "4321", nil },
			},
			lat: nil,
			ref: "4321",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			ref, err := Search(&fakeResolver{0, tc.com, tc.lat}, tc.env)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}

			if ref != tc.ref {
				t.Fatal("expected", tc.ref, "got", ref)
			}
		})
	}
}

// fakeResolver provides a controllable implementation of Resolver.Exists and
// Resolver.Latest, so that resolver.Search can be tested in isolation.
type fakeResolver struct {
	cal int
	com []func(string) (string, error)
	lat func() (string, error)
}

func (f *fakeResolver) Commit(ref string) (string, error) {
	defer func() { f.cal++ }()
	return f.com[f.cal](ref)
}

func (f *fakeResolver) Latest() (string, error) {
	return f.lat()
}
