package operator

import "github.com/xh3b4sd/tracer"

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

	// Lookup the artifact references for every service deployment that defines a
	// branch deployment strategy. This second step populates the desired artifact
	// reference.

	{
		err = h.ref.Ensure()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// TODO
	//
	//     get current state for docker image tags and ecs service tags
	//     get desired state for cloudformation templates from infra repo
	//     check for drift
	//     apply update, if any
	//     emit deployment event, if updated
	//

	return nil
}
