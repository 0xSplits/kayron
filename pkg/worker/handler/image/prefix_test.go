package image

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/xh3b4sd/logger"
)

func Test_Worker_Handler_Image_witPre_Sha(t *testing.T) {
	testCases := []struct {
		det []types.ImageDetail
		pre []types.ImageDetail
	}{
		// Case 000
		{
			det: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"v0.2.2"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"v1.0.0"}},
			},
			pre: nil,
		},
		// Case 001
		{
			det: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"3e935561e243157050afaab6a859ebf25ab7ca30"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"0458a0d88c600432ace658aa21ed357f0d24f2fc"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"d595ca48a2612597cc6e65cb104a66531edbefdd"}},
			},
			pre: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"3e935561e243157050afaab6a859ebf25ab7ca30"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"0458a0d88c600432ace658aa21ed357f0d24f2fc"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"d595ca48a2612597cc6e65cb104a66531edbefdd"}},
			},
		},
		// Case 003
		{
			det: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"v0.2.2"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"0458a0d88c600432ace658aa21ed357f0d24f2fc"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{""}},
			},
			pre: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"0458a0d88c600432ace658aa21ed357f0d24f2fc"}},
			},
		},
		// Case 003
		{
			det: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"x.y.z"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"sha256"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{""}},
			},
			pre: nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var han *Handler
			{
				han = &Handler{
					log: logger.Fake(),
				}
			}

			var opt []cmp.Option
			{
				opt = []cmp.Option{
					cmpopts.IgnoreUnexported(types.ImageDetail{}),
				}
			}

			pre := han.witPre(tc.det, strings.Split(string(Sha), ","))
			if dif := cmp.Diff(tc.pre, pre, opt...); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}

func Test_Worker_Handler_Image_witPre_Tag(t *testing.T) {
	testCases := []struct {
		det []types.ImageDetail
		pre []types.ImageDetail
	}{
		// Case 000
		{
			det: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"v0.2.2"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"v1.0.0"}},
			},
			pre: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"v0.2.2"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"v1.0.0"}},
			},
		},
		// Case 001
		{
			det: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"3e935561e243157050afaab6a859ebf25ab7ca30"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"0458a0d88c600432ace658aa21ed357f0d24f2fc"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"d595ca48a2612597cc6e65cb104a66531edbefdd"}},
			},
			pre: nil,
		},
		// Case 003
		{
			det: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"v0.2.2"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"0458a0d88c600432ace658aa21ed357f0d24f2fc"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{""}},
			},
			pre: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"v0.2.2"}},
			},
		},
		// Case 003
		{
			det: []types.ImageDetail{
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"x.y.z"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{"sha256"}},
				{ImagePushedAt: aws.Time(time.Unix(3, 0)), ImageTags: []string{""}},
			},
			pre: nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var han *Handler
			{
				han = &Handler{
					log: logger.Fake(),
				}
			}

			var opt []cmp.Option
			{
				opt = []cmp.Option{
					cmpopts.IgnoreUnexported(types.ImageDetail{}),
				}
			}

			pre := han.witPre(tc.det, strings.Split(string(Tag), ","))
			if dif := cmp.Diff(tc.pre, pre, opt...); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
