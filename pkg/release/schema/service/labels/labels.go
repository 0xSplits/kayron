package labels

// Struct contains runtime specific internals annotated inside the schema
// loader.
type Struct struct {
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
