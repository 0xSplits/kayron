package hash

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Hash struct {
	Dsh []byte
	Hsh []byte
}

func New(str string) Hash {
	sum := sha256.Sum256([]byte(str))
	enc := hex.EncodeToString(sum[:])
	cas := cases.Upper(language.English).String(enc[:8])

	return Hash{
		Dsh: []byte("-" + cas),
		Hsh: []byte(cas),
	}
}

func (h Hash) Empty() bool {
	return h.Dsh == nil && h.Hsh == nil
}

func (h Hash) String() string {
	return string(h.Hsh)
}
