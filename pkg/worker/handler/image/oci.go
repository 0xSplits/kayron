package image

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

// digest returns all of the existing container image digest hashes associated
// to the given list of image details, by resolving the digest groups according
// to every tagged image manifest. If an image tagged with v0.2.2 was built and
// pushed for multiple architectures, then we are looking for the image digest
// of the tagged image itself, plus for all image digests of the tagged image's
// children. E.g. https://github.com/docker/build-push-action produces triplets
// of images for every tagged image pushed by default.
func (h *Handler) digest(drp []types.ImageDetail) ([]string, error) {
	var dig [][]string
	{
		dig = make([][]string, len(drp))
	}

	fnc := func(i int, d types.ImageDetail) error {
		var err error

		var grp []string
		{
			grp, err = h.digGrp(h.imaTag(d))
			if err != nil {
				return tracer.Mask(err)
			}
		}

		// Assigning the digest group concurrently works because every parallel
		// function uses their own unique index on the collection slice that matches
		// the amount of work being done here.

		{
			dig[i] = grp
		}

		return nil
	}

	{
		err := parallel.Slice(drp, fnc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	// Just flatten all group digests into a flat list of digest hashes to have a
	// simple return value.

	var flt []string
	for _, x := range dig {
		flt = append(flt, x...)
	}

	return flt, nil
}

func (h *Handler) digGrp(tag string) ([]string, error) {
	var err error

	var ref name.Reference
	{
		ref, err = name.ParseReference(tag)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	h.log.Log(
		"level", "info",
		"message", "resolving digest group",
		"tag", tag,
		"ref", ref.Name(),
	)

	var des *remote.Descriptor
	{
		des, err = remote.Get(ref, remote.WithAuthFromKeychain(h.key))
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	// Try to resolve the digest group of the given image tag based on the
	// underlying media type. Our expectation is to work with tagged image
	// indices, because we push multi architecture images. The exception should be
	// singular images that have been built and pushed in a non standard way,
	// maybe even by accident. If the underlying media type is neither an index
	// nor an image, then we return an error in order for a human to investigate
	// further.

	var dig []string

	if des.MediaType.IsImage() {
		dig = []string{
			des.Digest.String(),
		}
	}

	if des.MediaType.IsIndex() {
		dig, err = h.digInd(des)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	if len(dig) == 0 {
		return nil, tracer.Mask(invalidImageTagError,
			tracer.Context{Key: "media-type", Value: string(des.MediaType)},
			tracer.Context{Key: "tag", Value: tag},
		)
	} else {
		h.log.Log(
			"level", "info",
			"message", "resolved digest group",
			"tag", tag,
			"digest", strings.Join(dig, ","),
		)
	}

	return dig, nil
}

func (h *Handler) digInd(des *remote.Descriptor) ([]string, error) {
	var err error

	var ind v1.ImageIndex
	{
		ind, err = des.ImageIndex()
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var man *v1.IndexManifest
	{
		man, err = ind.IndexManifest()
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	// Note that we should put the tagged index digest last, because this parent
	// hash is the parent identifier of the entire digest group by virtue of its
	// index manifest. If the deletion process ever changes, then we ensure that
	// we delete all the children first, so that during failover we can retry and
	// discover the remain digest group over and over again, until the parent
	// itself is gone.

	var dig []string
	for _, x := range man.Manifests {
		dig = append(dig, x.Digest.String())
	}

	{
		dig = append(dig, des.Digest.String())
	}

	return dig, nil
}

func (h *Handler) imaTag(d types.ImageDetail) string {
	return fmt.Sprintf(
		"%s.dkr.ecr.%s.amazonaws.com/%s:%s", // e.g. 995626699990.dkr.ecr.us-west-2.amazonaws.com/kayron:v0.2.2
		ptrStr(d.RegistryId),
		h.ecr.Options().Region,
		ptrStr(d.RepositoryName),
		slcStr(d.ImageTags), // h.witPre must guarantee exactly 1 image tag
	)
}

func ptrStr(ptr *string) string {
	if ptr != nil {
		return *ptr
	}

	return ""
}

func slcStr(slc []string) string {
	if len(slc) != 0 {
		return slc[0]
	}

	return ""
}
