package preview

import (
	"bytes"
)

func appTag(lin []byte, hsh string) ([]byte, error) {
	var pre []byte
	{
		pre = bytes.Join(
			[][]byte{
				nil,
				[]byte("        - Key: \"preview\""),
				[]byte("          Value: \"" + hsh + "\""),
			},
			[]byte("\n"),
		)
	}

	var byt []byte
	{
		byt = append(lin, pre...)
	}

	return byt, nil
}
