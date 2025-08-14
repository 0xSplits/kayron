package stack

import "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

type Interface interface {
	// Delete purges the underlying local cache, causing Search to fetch the
	// latest version of the stack object state again over network.
	Delete()

	// Search returns the state of the configured stack object and caches the
	// first valid search result so that consecutive executions of Search prevent
	// network calls. This behaviour guarantees consistent stack object state
	// within reconciliation loops.
	Search() (types.Stack, error)
}
