package roghfs

import (
	"os"

	"github.com/xh3b4sd/tracer"
)

// Stat tries to return an instance of os.FileInfo for the given file path. Note
// that Stat() is called before every loop of the walk function of afero.Walk(),
// because Stat() provides the fs.FileInfo instance for every walk function
// call. If the given file does not exist, an os.PathError is returned.
func (r *Roghfs) Stat(pat string) (os.FileInfo, error) {
	var err error

	var fil os.FileInfo
	{
		fil, err = r.bas.Stat(pat)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return fil, nil
}
