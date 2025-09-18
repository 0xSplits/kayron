package labels

import (
	"github.com/0xSplits/kayron/pkg/hash"
)

// Struct contains runtime specific internals annotated inside the schema
// loader.
type Struct struct {
	// Block indicates the index of the service definition relative to the context
	// of a specific config file. E.g. given a config file with 3 service
	// definitions, the last block has the index 2.
	Block int
	// Hash contains the hashed branch name for any service release of a preview
	// deployment.
	Hash hash.Hash
	// Head is the latest Git Reference for any service release of a preview
	// deployment.
	Head string
	// Source is the absolute source file path of the .yaml definition as loaded
	// from the underlying file system. This label may help to make error messages
	// more useful.
	Source string
}

func (s Struct) Empty() bool {
	return s.Block == 0 && s.Hash.Empty() && s.Source == ""
}

func (s Struct) Verify() error {
	return nil
}
