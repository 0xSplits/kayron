package registry

import (
	"github.com/xh3b4sd/tracer"
)

var metricNameWhitelistError = &tracer.Error{
	Kind: "metricNameWhitelistError",
	Desc: "The caller used a metric name that is not registered in the whitelist.",
}

var labelKeyWhitelistError = &tracer.Error{
	Kind: "labelKeyWhitelistError",
	Desc: "The caller used a label key that is not registered in the whitelist.",
}

var labelValueWhitelistError = &tracer.Error{
	Kind: "labelValueWhitelistError",
	Desc: "The caller used a label value that is not registered in the whitelist.",
}
