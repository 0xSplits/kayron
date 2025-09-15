package hash

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Hash struct {
	Hsh []byte
	Dsh []byte
}

func New(str string) Hash {
	sum := sha256.Sum256([]byte(str))
	enc := hex.EncodeToString(sum[:])
	cas := cases.Upper(language.English).String(enc[:8])

	return Hash{
		Hsh: []byte(cas),
		Dsh: []byte("-" + cas),
	}
}
