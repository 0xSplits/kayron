package webhook

import (
	"github.com/xh3b4sd/tracer"
)

var commitHashEmptyError = &tracer.Error{
	Description: "This error indicates that no hash could be found for the commit of the push event, which means that the webhook cache can neither identify nor store this specific webhook payload.",
}

var commitTimeEmptyError = &tracer.Error{
	Description: "This error indicates that no timestamp could be found for the commit of the push event, which means that the webhook cache can neither identify nor store this specific webhook payload.",
}
