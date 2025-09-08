package image

import (
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/xh3b4sd/tracer"
)

// Ensure implements the following algorithm to manage registry bloat.
//
//	K    threshold amount of tagged images to keep
//	D    threshold amount of tagged images to drop
//
//	fetch e.g. 100 tagged images
//
//	create a list with tag prefixes, e.g. sha or tag
//
//	if <= K, break
//	sort by date, oldest first
//	create list of up to D items, never going beyond K
//	for each D item
//	    fetch digests
//
//	batch delete collected digests
func (h *Handler) Ensure() error {
	var err error

	// At first, fetch a sufficiently large list of container images for the
	// configured repository. The amount of images we are looking for here must be
	// larger than the amount of most recent image tags that we want to keep in
	// place. If we want to keep 10 of the most recent image tags, then we have to
	// fetch more than that. The default of images to lookup here is 100.

	var ima []types.ImageDetail
	{
		ima, err = h.search(string(h.rep))
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// From the entire list of images fetched above, we filter out those that do
	// not match our handler configured image tag prefixes, leaving us with a list
	// of image details matching our desired set of tag prefixes.

	var pre []types.ImageDetail
	{
		pre = h.witPre(ima, h.pre)
	}

	// If there are less images in the registry than we actually want to maintain
	// at all times, then we stop processing here.

	if len(pre) <= Keep {
		return nil
	}

	// TODO

	var rem []types.ImageDetail
	{
		rem = selRem(pre)
	}

	// TODO

	var dig []string
	{
		dig, err = h.digest(rem)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// TODO

	{
		err = h.delete(string(h.rep), dig)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
