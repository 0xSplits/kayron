package image

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_Worker_Handler_Image_selRem(t *testing.T) {
	testCases := []struct {
		pre []types.ImageDetail
		rem []types.ImageDetail
	}{
		// Case 000, empty
		{
			pre: nil,
			rem: nil,
		},
		// Case 001, keep all
		{
			pre: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(1, 0)), ImageTags: []string{"img-a"}},
			},
			rem: nil,
		},
		// Case 002, keep all
		{
			pre: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(1, 0)), ImageTags: []string{"img-a"}},
				{ImagePushedAt: aws.Time(time.Unix(2, 0)), ImageTags: []string{"img-b"}},
			},
			rem: nil,
		},
		// Case 003, remove oldest
		{
			pre: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"newest"}},
				{ImagePushedAt: aws.Time(time.Unix(1, 0)), ImageTags: []string{"oldest"}},
				{ImagePushedAt: aws.Time(time.Unix(2, 0)), ImageTags: []string{"middle"}},
			},
			rem: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(1, 0)), ImageTags: []string{"oldest"}},
			},
		},
		// Case 004, remove max
		{
			pre: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"t3"}},
				{ImagePushedAt: aws.Time(time.Unix(4, 0)), ImageTags: []string{"t4"}},
				{ImagePushedAt: aws.Time(time.Unix(1, 0)), ImageTags: []string{"t1"}},
				{ImagePushedAt: aws.Time(time.Unix(6, 0)), ImageTags: []string{"t6"}},
				{ImagePushedAt: aws.Time(time.Unix(5, 0)), ImageTags: []string{"t5"}},
				{ImagePushedAt: aws.Time(time.Unix(2, 0)), ImageTags: []string{"t2"}},
			},
			rem: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(1, 0)), ImageTags: []string{"t1"}},
				{ImagePushedAt: aws.Time(time.Unix(2, 0)), ImageTags: []string{"t2"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"t3"}},
			},
		},
		// Case 005, remove and keep exact, no panic
		{
			pre: []types.ImageDetail{
				{ /**************** nil *****************/ ImageTags: []string{"t0"}},
				{ImagePushedAt: aws.Time(time.Unix(4, 0)), ImageTags: []string{"t4"}},
				{ImagePushedAt: aws.Time(time.Unix(1, 0)), ImageTags: []string{"t1"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"t3"}},
				{ImagePushedAt: aws.Time(time.Unix(2, 0)), ImageTags: []string{"t2"}},
			},
			rem: []types.ImageDetail{
				{ /**************** nil *****************/ ImageTags: []string{"t0"}},
				{ImagePushedAt: aws.Time(time.Unix(1, 0)), ImageTags: []string{"t1"}},
				{ImagePushedAt: aws.Time(time.Unix(2, 0)), ImageTags: []string{"t2"}},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var rem []types.ImageDetail
			{
				rem = selRem(tc.pre, 3, 2) // drop 3 at a time, keep 2 in the list
			}

			{
				opt := []cmp.Option{
					cmpopts.IgnoreUnexported(types.ImageDetail{}),
				}

				if dif := cmp.Diff(tc.rem, rem, opt...); dif != "" {
					t.Fatalf("-expected +actual:\n%s", dif)
				}
			}
		})
	}
}
