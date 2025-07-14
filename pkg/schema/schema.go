// Package schema defines the in-memory representation of Kayron’s
// change-management YAML and carries the field-level tags expected by
// goccy/go-yaml.
package schema

import (
	"github.com/0xSplits/kayron/pkg/schema/infrastructure"
	service "github.com/0xSplits/kayron/pkg/schema/service"
	yaml "github.com/goccy/go-yaml"
)

var _ = yaml.Marshal // TODO

// Schema is the root of our deployment configuration.
type Schema struct {
	Infrastructure infrastructure.Infrastructure `yaml:"infrastructure,omitempty"`
	Service        []service.Service             `yaml:"service,omitempty"`
}

// TODO validate unique values across multiple environments, e.g. must not use
// prod twice as env shorthand
