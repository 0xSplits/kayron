package image

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

type prefix string

const (
	// Sha represents the set of prefixes matching image tags in the form of
	// hexadecimal commit hashes, e.g. 1cd6...6863.
	Sha prefix = "0,1,2,3,4,5,6,7,8,9,a,b,c,d,e,f"

	// Tag represents the prefix matching image tags in the form of semver
	// versioned release tags, e.g. v0.2.2.
	Tag prefix = "v"
)

// TODO write tests
func (h *Handler) witPre(lis []types.ImageDetail, pre []string) []types.ImageDetail {
	var wit []types.ImageDetail

	for _, x := range lis {
		// The image detail response must meet certain requirements so that our
		// algorithm can work as expected.
		//
		//     1. We cannot work with nil values for push times, because our
		//        algorithm relies on a list of image tags, sorted by age, from
		//        oldest to newest.
		//
		//     2. We cannot work with images that do not have exactly 1 image tag,
		//        because that image tag must either represent a container image
		//        of a commit sha or release tag.
		//

		if x.ImagePushedAt == nil || len(x.ImageTags) != 1 {
			h.log.Log(
				"level", "warning",
				"message", "skipping cleanup for image tag",
				"reason", "invalid image push timestamp or image tag",
				"digest", *x.ImageDigest,
				"registry", *x.RegistryId,
				"repository", *x.RepositoryName,
				"tags", strings.Join(x.ImageTags, ","),
			)

			{
				continue
			}
		}

		if hasPre(x.ImageTags, pre) {
			wit = append(wit, x)
		}
	}

	return wit
}

func hasPre(lis []string, pre []string) bool {
	for _, x := range lis {
		for _, y := range pre {
			if strings.HasPrefix(x, y) {
				return true
			}
		}
	}

	return false
}
