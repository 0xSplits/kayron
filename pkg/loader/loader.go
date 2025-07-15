package loader

import (
	"io/fs"
	"path/filepath"

	"github.com/0xSplits/kayron/pkg/schema"
	yaml "github.com/goccy/go-yaml"
	"github.com/spf13/afero"
)

// Loader walks the provided filesystem starting at the given root directory.
// All .yaml files are inspected and their content is marshalled into Schema
// types. An error is returned if walking, reading or unmarshalling fails.
func Loader(fil afero.Fs, roo string) (schema.Schemas, error) {
	var lis schema.Schemas

	wal := func(pat string, inf fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if inf.IsDir() {
			return nil
		}
		if filepath.Ext(inf.Name()) != ".yaml" {
			return nil
		}

		raw, err := afero.ReadFile(fil, pat)
		if err != nil {
			return err
		}

		var sch schema.Schema
		{
			err := yaml.Unmarshal(raw, &sch)
			if err != nil {
				return err
			}
		}

		{
			lis = append(lis, sch)
		}

		return nil
	}

	err := afero.Walk(fil, filepath.Clean(roo), wal)
	if err != nil {
		return nil, err
	}

	return lis, nil
}
