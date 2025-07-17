package specification

import (
	"errors"

	"github.com/xh3b4sd/tracer"
)

var schemaEmptyError = &tracer.Error{
	Kind: "schemaEmptyError",
	Desc: "The schema configuration requires at least one service to be provided.",
}

func IsSchemaEmpty(err error) bool {
	return errors.Is(err, schemaEmptyError)
}

var infrastructureShorthandError = &tracer.Error{
	Kind: "infrastructureShorthandError",
	Desc: "The infrastructure configuration requires all shorthands to be unique.",
}

func IsInfrastructureShorthand(err error) bool {
	return errors.Is(err, infrastructureShorthandError)
}
