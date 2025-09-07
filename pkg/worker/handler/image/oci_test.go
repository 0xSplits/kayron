package image

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/google/go-cmp/cmp"
)

func Test_Worker_Handler_Image_imaTag(t *testing.T) {
	testCases := []struct {
		det types.ImageDetail
		reg string
		tag string
	}{
		// Case 000
		{
			det: types.ImageDetail{
				RegistryId:     aws.String("995626699990"),
				RepositoryName: aws.String("kayron"),
				ImageTags:      []string{"v0.2.2"},
			},
			reg: "us-west-2",
			tag: "995626699990.dkr.ecr.us-west-2.amazonaws.com/kayron:v0.2.2",
		},
		// Case 000
		{
			det: types.ImageDetail{
				RegistryId:     aws.String("1234"),
				RepositoryName: aws.String("specta"),
				ImageTags:      []string{"v1.0.0"},
			},
			reg: "eu-central-1",
			tag: "1234.dkr.ecr.eu-central-1.amazonaws.com/specta:v1.0.0",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var han *Handler
			{
				han = &Handler{
					ecr: ecr.NewFromConfig(aws.Config{
						Region: tc.reg,
					}),
				}
			}

			tag := han.imaTag(tc.det)
			if dif := cmp.Diff(tc.tag, tag); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
