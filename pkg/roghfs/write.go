package roghfs

import (
	"os"

	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"
)

// write is an internal helper to persist the given source file bytes in the
// underlying base file system using the provided file path.
func (r *Roghfs) write(pat string, byt []byte) error {
	var err error

	var fil afero.File
	{
		fil, err = r.bas.OpenFile(pat, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o444)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		defer fil.Close() // nolint:errcheck
	}

	{
		_, err = fil.Write(byt)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
