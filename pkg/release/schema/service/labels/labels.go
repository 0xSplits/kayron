package labels

// Struct contains runtime specific internals annotated inside the schema
// loader.
type Struct struct {
	// Block indicates the index of the service definition relative to the context
	// of a specific config file. E.g. given a config file with 3 service
	// definitions, the last block has the index 2.
	Block int
	// Source is the absolute source file path of the .yaml definition as loaded
	// from the underlying file system. This label may help to make error messages
	// more useful.
	Source string
}

func (m Struct) Empty() bool {
	return m.Source == ""
}

func (m Struct) Verify() error {
	return nil
}
