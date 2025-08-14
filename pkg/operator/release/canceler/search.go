package canceler

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
		types.StackStatusUpdateRollbackFailed: true,

		// healthy final states
		types.StackStatusCreateComplete:         false,
		types.StackStatusImportComplete:         false,
		types.StackStatusImportRollbackComplete: false,
		types.StackStatusUpdateComplete:         false,
		types.StackStatusUpdateRollbackComplete: false,
	}
)

func (c *Canceler) Cancel() (bool, error) {
	var err error

	var sta types.Stack
	{
		sta, err = c.sta.Search()
		if err != nil {
			return false, tracer.Mask(err)
		}
	}

	var can bool
	var exi bool
	{
		can, exi = status[sta.StackStatus]
	}

	if !exi {
		return false, tracer.Mask(invalidStackStatusError,
			tracer.Context{Key: "stack", Value: c.env.CloudformationStack},
			tracer.Context{Key: "status", Value: sta},
		)
	}

	return can, nil
}
