package roghfs

import (
	"testing"

	"github.com/spf13/afero"
)

// Test_Roghfs_Interface ensures that this read-only file system implementation
// complies with the afero file system interface. This test already fails at
// compile time if Roghfs does not implement afero.Fs.
func Test_Roghfs_Interface(t *testing.T) {
	var _ afero.Fs = &Roghfs{}
}
