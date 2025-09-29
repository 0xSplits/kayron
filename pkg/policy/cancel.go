package policy

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/xh3b4sd/tracer"
)

var (
	// status is a mapping of stack statuses answering the question if the current
	// reconciliation loop should be cancelled given a certain stack status. E.g.
	// should we cancel on UPDATE_COMPLETE=false? No, because the underlying
	// CloudFormation stack is ready to be reconciled again after it has been
	// updated previously.
	status = map[types.StackStatus]bool{
		// transitional progress states
		types.StackStatusCreateInProgress:                        true,
		types.StackStatusDeleteInProgress:                        true,
		types.StackStatusImportInProgress:                        true,
		types.StackStatusReviewInProgress:                        true,
		types.StackStatusUpdateCompleteCleanupInProgress:         true,
		types.StackStatusUpdateInProgress:                        true,
		types.StackStatusImportRollbackInProgress:                true,
		types.StackStatusRollbackInProgress:                      true,
		types.StackStatusUpdateRollbackCompleteCleanupInProgress: true,
		types.StackStatusUpdateRollbackInProgress:                true,

		// unrecoverable final states
		types.StackStatusCreateFailed:         true,
		types.StackStatusDeleteFailed:         true,
		types.StackStatusDeleteComplete:       true,
		types.StackStatusImportRollbackFailed: true,
		types.StackStatusRollbackComplete:     true,
		types.StackStatusRollbackFailed:       true,
		types.StackStatusUpdateFailed:         true,

		// healthy final states
		types.StackStatusUpdateRollbackFailed:   false, // TODO move up again
		types.StackStatusCreateComplete:         false,
		types.StackStatusImportComplete:         false,
		types.StackStatusImportRollbackComplete: false,
		types.StackStatusUpdateComplete:         false,
		types.StackStatusUpdateRollbackComplete: false,
	}
)

// Cancel tells us whether it is safe to proceed with the next CloudFormation
// update this time around. Note that Cancel uses the stack cache that is being
// purged at the start of every reconciliation loop. Further note that any
// internal error causes Cancel to return true, which is meant to stop
// processing in case our understanding of the current state of the system is
// incomplete.
func (p *Policy) Cancel() bool {
	var err error

	var sta types.Stack
	{
		sta, err = p.Stack()
		if err != nil {
			p.log.Log(
				"level", "error",
				"message", "stack search error",
				"stack", tracer.Json(err),
			)

			return true
		}
	}

	var can bool
	var exi bool
	{
		can, exi = status[sta.StackStatus]
	}

	if !exi {
		p.log.Log(
			"level", "error",
			"message", "stack search error",
			"stack", tracer.Json(invalidStackStatusError),
		)

		return true
	}

	return can
}
