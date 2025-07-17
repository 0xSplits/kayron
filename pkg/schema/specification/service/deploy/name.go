package deploy

import (
	"fmt"
	"strings"
)

// name returns the package declaration of the given Interface implementation.
func name(h Interface) string {
	//
	//     *release.Release
	//
	var p string
	{
		p = fmt.Sprintf("%T", h)
	}

	//
	//     release.Release
	//
	var t string
	{
		t = strings.TrimPrefix(p, "*")
	}

	//
	//     release
	//
	var s string
	{
		s = strings.Split(t, ".")[0]
	}

	return s
}
