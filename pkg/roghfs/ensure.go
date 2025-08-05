package roghfs

import (
	"os"
	"path/filepath"

	"github.com/xh3b4sd/tracer"
)

// ensure guarantees that the given file or directory exists inside the injected
// base file system. If the requested file path is not cached locally, then
// ensure fetches its raw bytes from GitHub.
func (r *Roghfs) ensure(pat string) error {
	var err error

	// Ensure a clean and standardized file path string.
	{
		pat = filepath.Clean(pat)
	}

	// If we already fetched the content of the requested source file, then we
	// don't have to do any more work here, because ensure() did its job already
	// in a previous call.
	{
		exi := r.cac.Exists(pat)
		if exi {
			return nil
		}
	}

	var fil os.FileInfo
	{
		fil, err = r.bas.Stat(pat)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// If the requested path points to a directory, then we do not have to fetch
	// any more content, because we already initialized the entirety of the file
	// system tree.
	{
		dir := fil.IsDir()
		if dir {
			return nil
		}
	}

	// Download the requested source file content from the remote Github
	// repository.
	var byt []byte
	{
		byt, err = r.search(pat)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Persist the downloaded source file content in the underlying base file
	// system.
	{
		err = r.write(pat, byt)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Finally mark the requested source file as being cached, and surface all
	// errors related to our internal cache logic, if any.
	{
		exi := r.cac.Create(pat, struct{}{})
		if exi {
			return tracer.Mask(fileAlreadyCachedError, tracer.Context{Key: "path", Value: pat})
		}
	}

	return nil
}
