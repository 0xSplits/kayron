package webhook

import (
	"fmt"
	"strings"
)

// Slice contains a list of alternative deployment mechanisms. Each webhook
// provided here is invoked to deploy e.g. our frontends in Vercel. The format
// of those webhook definitions requires the usage of a prefix for the HTTP
// method that this webhook should be called with. It is further required to
// provide a HTTPs URL. Failed webhook calls may be retried and eventually be
// reported as terminal failure.
//
//	POST:https://{{DNS}}/{{PATH}}
type Slice []String

func (s Slice) Empty() bool {
	return len(s) == 0
}

func (s Slice) String() string {
	if len(s) == 0 {
		return "[]"
	}

	var lis []string
	for _, x := range s {
		lis = append(lis, x.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(lis, ","))
}

func (s Slice) Verify() error {
	// TODO
	return nil
}
