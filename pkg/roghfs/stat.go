package roghfs

import (
	"os"

	"github.com/xh3b4sd/tracer"
)

// Stat tries to return an instance of os.FileInfo for the given file path. Note
// that Stat is called before every loop of the walk function of afero.Walk,
// because Stat provides the fs.FileInfo instance for every walk function call.
// The first call of Stat will also setup the entire file and folder structure
// within the underlying base file system. If the given file does not exist
// after all, an os.PathError is returned.
func (r *Roghfs) Stat(pat string) (os.FileInfo, error) {
	var err error

	// Bootstrap the entire file system tree according to the configured remote
	// Github repository, but only if we have not initialized the base file system
	// successfully before. This ensures that we traverse the entire repository
	// tree only one time, while only making a single network call. Note that this
	// must be done here in Stat, because an empty base file system cannot walk a
	// specific remote file path, which is why we have to bootstrap the file
	// system tree before calling Stat on the base file system.

	{
		err = r.mut.Success(r.tree)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var fil os.FileInfo
	{
		fil, err = r.bas.Stat(pat)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return fil, nil
}
