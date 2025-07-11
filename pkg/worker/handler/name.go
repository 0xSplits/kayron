package handler

import (
	"fmt"
	"strings"
)

// Name returns the package declaration of the given handler implementation.
func Name(h Interface) string {
	//
	//     *keypair.Handler
	//
	var p string
	{
		p = fmt.Sprintf("%T", h)
	}

	//
	//     keypair.Handler
	//
	var t string
	{
		t = strings.TrimPrefix(p, "*")
	}

	//
	//     keypair
	//
	var s string
	{
		s = strings.Split(t, ".")[0]
	}

	return s
}
