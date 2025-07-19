package loader

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/0xSplits/kayron/pkg/schema/specification"
	yaml "github.com/goccy/go-yaml"
	"github.com/spf13/afero"
)

// Loader walks the provided filesystem starting at the given root directory.
// All .yaml files are inspected and their content is marshalled into Schema
// types. An error is returned if walking, reading or unmarshalling fails.
func Loader(sys afero.Fs, roo string) (specification.Schemas, error) {
	var lis specification.Schemas

	wal := func(pat string, fil fs.FileInfo, err error) error {
		{
			if err != nil {
				return err
			}
			if fil.IsDir() {
				return nil
			}
		}

		var ext string
		{
			ext = filepath.Ext(fil.Name())
			if ext != ".yaml" {
				return nil
			}
		}

		var byt []byte
		{
			byt, err = afero.ReadFile(sys, pat)
			if err != nil {
				return err
			}
		}

		var sch specification.Schema
		{
			err := yaml.Unmarshal(byt, &sch)
			if err != nil {
				return err
			}
		}

		{
			sch.Labels.Environment = strings.TrimSuffix(fil.Name(), ext)
			sch.Labels.Source = pat
			sch.Labels.Testing = filepath.Base(filepath.Dir(pat)) == "testing"
		}

		{
			lis = append(lis, sch)
		}

		return nil
	}

	err := afero.Walk(sys, filepath.Clean(roo), wal)
	if err != nil {
		return nil, err
	}

	return lis, nil
}
