package image

import (
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

// selRem returns those image details that are supposed to be cleaned up, based
// on the given drop and keep parameters. The returned list of image details is
// either nil, or sorted by push timestamp.
func selRem(pre []types.ImageDetail, dro int, kee int) []types.ImageDetail {
	// sort by oldest first, h.witPre must guarantee a non nil push timestamp

	sort.Slice(pre, func(i, j int) bool {
		return aws.ToTime(pre[i].ImagePushedAt).Before(aws.ToTime(pre[j].ImagePushedAt))
	})

	var num int
	{
		num = min(dro, max(0, len(pre)-kee))
	}

	// If there is nothing to do, return nil early.

	if num == 0 {
		return nil
	}

	var rem []types.ImageDetail
	{
		rem = pre[:num]
	}

	return rem
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
