package preview

import (
	"fmt"
	"strings"
)

func preDom(hsh string, doc string) string {
	// Note that this is a dirty hack to make preview deployments work today for
	// existing services that already work using certain incosnistencies between
	// repository and domain names. E.g. we have "splits-lite" in Github, but use
	// just "lite.testing.splits.org". A better way of doing this would be to allow
	// for some kind of domain configuration in the release definition, so that we
	// can remove this magical string replacement below.

	var trm string
	{
		trm = strings.TrimPrefix(doc, "splits-")
	}

	return fmt.Sprintf("%s.%s.${Environment}.splits.org", hsh, trm)
}
