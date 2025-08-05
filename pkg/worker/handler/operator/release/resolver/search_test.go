package resolver

import (
	"fmt"
	"testing"
)

func Test_Worker_Handler_Releases_Resolver_Search_failure(t *testing.T) {
	testCases := []struct {
		env string
		com func() (string, error)
		exi func() (bool, error)
		lat func() (string, error)
		mat func(error) bool
	}{
		// Case 000, production, no branch, no release
		{
			env: "production",
			com: func() (string, error) { return "", nil },
			exi: func() (bool, error) { return false, nil },
			lat: func() (string, error) { return "", releaseNotFoundError },
			mat: IsReleaseNotFound,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			_, err := Search(fakeResolver{tc.com, tc.exi, tc.lat}, tc.env)
			if !tc.mat(err) {
				t.Fatal("expected", true, "got", err)
			}
		})
	}
}

func Test_Worker_Handler_Releases_Resolver_Search_success(t *testing.T) {
	testCases := []struct {
		env string
		com func() (string, error)
		exi func() (bool, error)
		lat func() (string, error)
		ref string
	}{
		// Case 000, staging, default branch
		{
			env: "staging",
			com: func() (string, error) { return "1234", nil },
			exi: func() (bool, error) { return false, nil },
			lat: func() (string, error) { return "", nil },
			ref: "1234",
		},
		// Case 001, production, no branch
		{
			env: "production",
			com: func() (string, error) { return "", nil },
			exi: func() (bool, error) { return false, nil },
			lat: func() (string, error) { return "v0.1.0", nil },
			ref: "v0.1.0",
		},
		// Case 002, production, no release
		{
			env: "production",
			com: func() (string, error) { return "", nil },
			exi: func() (bool, error) { return true, nil },
			lat: func() (string, error) { return "", releaseNotFoundError },
			ref: "production",
		},
		// Case 003, testing, no branch
		{
			env: "testing",
			com: func() (string, error) { return "1234", nil },
			exi: func() (bool, error) { return false, nil },
			lat: func() (string, error) { return "", releaseNotFoundError },
			ref: "1234",
		},
		// Case 004, arbitrary test environment
		{
			env: "melissa",
			com: func() (string, error) { return "", nil },
			exi: func() (bool, error) { return true, nil },
			lat: func() (string, error) { return "", releaseNotFoundError },
			ref: "melissa",
		},
		// Case 005, arbitrary test environment, no branch
		{
			env: "melissa",
			com: func() (string, error) { return "1234", nil },
			exi: func() (bool, error) { return false, nil },
			lat: func() (string, error) { return "", releaseNotFoundError },
			ref: "1234",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			ref, err := Search(fakeResolver{tc.com, tc.exi, tc.lat}, tc.env)
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
	com func() (string, error)
	exi func() (bool, error)
	lat func() (string, error)
}

func (f fakeResolver) Commit() (string, error) {
	return f.com()
}

func (f fakeResolver) Exists(bra string) (bool, error) {
	return f.exi()
}

func (f fakeResolver) Latest() (string, error) {
	return f.lat()
}
