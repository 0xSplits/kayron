//go:build integration

package roghfs

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-github/v73/github"
	"github.com/spf13/afero"
)

// Test_Roghfs_Integration runs the read‑only Github file system against a
// remote repository and compares 3 selected files as they were created in
// Kayron's first pull request against our local golden files. The required auth
// token to run this integration test needs at least public repo permissions.
//
//	KAYRON_GITHUB_TOKEN=todo go test -tags=integration ./...
func Test_Roghfs_Integration(t *testing.T) {
	var gfs *Roghfs
	{
		gfs = New(Config{
			Bas: afero.NewMemMapFs(),
			Git: github.NewClient(nil).WithAuthToken(musTok()),
			Org: "0xSplits",
			Rep: "kayron",
			Ref: "d2f2a18b998172039c6f2a325d4c83de20819e3e", // setup project structure for prototype (#1)
		})
	}

	var sys afero.Fs
	{
		sys = afero.NewReadOnlyFs(afero.NewOsFs())
	}

	var wht []string
	{
		wht = []string{
			".github/workflows/go-build.yaml",
			"pkg/worker/handler/operator/cooler.go",
			"go.mod",
		}
	}

	// Collect the local and remote data using the host's file system and our
	// read-only Github file system. Below we compare the contents of our
	// whitelisted file paths. If the raw bytes that our Github file system foud
	// remotely match our hard coded golden file copies, we know that our Github
	// file system implementation works properly.

	var loc [][]byte
	var rem [][]byte
	{
		loc = musWlk(sys, "testdata", wht)
		rem = musWlk(gfs, ".", wht)
	}

	{
		if len(loc) != len(wht) {
			t.Fatalf("expected %#v got %#v", len(wht), len(loc))
		}
		if len(rem) != len(wht) {
			t.Fatalf("expected %#v got %#v", len(wht), len(rem))
		}
	}

	for i := range wht {
		if len(loc[i]) == 0 {
			t.Fatal("expected the local file system to contain content bytes, got empty file")
		}
		if len(rem[i]) == 0 {
			t.Fatal("expected the remote file system to contain content bytes, got empty file")
		}
		if !bytes.Equal(loc[i], rem[i]) {
			t.Fatal("expected both file systems to contain equal bytes, got different files")
		}
	}
}

func musTok() string {
	var tok string
	{
		tok = os.Getenv("KAYRON_GITHUB_TOKEN")
		if tok == "" {
			panic("env var KAYRON_GITHUB_TOKEN must not be empty")
		}
	}

	return tok
}

// musWlk simulates a loader function using afero.Walk so that we can collect
// the file content of local and remote source files.
func musWlk(afs afero.Fs, roo string, wht []string) [][]byte {
	var lis [][]byte

	wlk := func(pat string, fil fs.FileInfo, err error) error {
		{
			if err != nil {
				panic(err)
			}
			if fil.IsDir() {
				return nil
			}
			if !slices.Contains(wht, trmPat(pat)) {
				return nil
			}
		}

		var byt []byte
		{
			byt, err = afero.ReadFile(afs, pat)
			if err != nil {
				panic(err)
			}
		}

		{
			lis = append(lis, byt)
		}

		return nil
	}

	{
		err := afero.Walk(afs, filepath.Clean(roo), wlk)
		if err != nil {
			panic(err)
		}
	}

	return lis
}

// trmPat normalizes the file paths of the local golden files so that we have
// one standard way of filtering local and remote files by file path.
func trmPat(pat string) string {
	pat = strings.TrimPrefix(pat, "testdata/")
	pat = strings.TrimSuffix(pat, ".golden")

	return pat
}
