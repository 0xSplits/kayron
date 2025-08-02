package operator

import (
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

func (h *Handler) Ensure() error {
	var err error

	// Lookup all the release settings for the configured releases repository.
	// This first step initializes the release and artifact caches.

	{
		err = h.rel.Ensure()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Run the next steps in parallel in order to find the current and desired
	// state of the service releases that we are tasked to managed.
	//
	//     1. Lookup the ECS container of every service release regardless of its
	//        deployment strategy. This populates the CURRENT state of the
	//        artifact reference.
	//
	//     2. Lookup the Github reference of every service release that defines a
	//        branch deployment strategy. This populates the DESIRED state of the
	//        artifact reference.
	//
	//     3. TODO fetch existing image tags from configured Docker registry
	//
	//     4. TODO fetch existing cloudformation templates from infrastructure repo
	//

	{
		err = parallel.Func(h.con.Ensure, h.ref.Ensure)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// TODO
	//
	//     check for drift
	//     apply update, if any
	//     emit deployment event, if updated
	//

	return nil
}
