package image

import (
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

// TODO write tests
func selRem(pre []types.ImageDetail) []types.ImageDetail {
	// sort by oldest first
	// h.witPre must guarantee a non nil push timestamp

	sort.Slice(pre, func(i, j int) bool {
		return pre[i].ImagePushedAt.Before(*pre[j].ImagePushedAt)
	})

	var num int
	{
		num = min(Drop, max(0, len(pre)-Keep))
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
