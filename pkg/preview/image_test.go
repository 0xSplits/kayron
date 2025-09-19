package preview

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Preview_repIma(t *testing.T) {
	testCases := []struct {
		lin []byte
		tag []byte
		out []byte
	}{
		// Case 000
		{
			lin: []byte(`Image: !Sub "${AWS::AccountId}.dkr.ecr.${AWS::Region}.amazonaws.com/splits-lite:${LiteVersion}"`),
			tag: []byte(`v0.2.0`),
			out: []byte(`Image: !Sub "${AWS::AccountId}.dkr.ecr.${AWS::Region}.amazonaws.com/splits-lite:v0.2.0"`),
		},
		// Case 001
		{
			lin: []byte(`          Image: 'ecr.amazonaws.com/splits-lite:${LiteVersion}' # hello world`),
			tag: []byte(`v1.0.0`),
			out: []byte(`          Image: 'ecr.amazonaws.com/splits-lite:v1.0.0' # hello world`),
		},
		// Case 002
		{
			lin: []byte(`    Image: !Sub "${AWS::AccountId}.dkr.ecr.${AWS::Region}.amazonaws.com/splits-lite:v0.1.0"`),
			tag: []byte(`bc7891268e44f62e0aebbe339c0850b61d52c417`),
			out: []byte(`    Image: !Sub "${AWS::AccountId}.dkr.ecr.${AWS::Region}.amazonaws.com/splits-lite:bc7891268e44f62e0aebbe339c0850b61d52c417"`),
		},
		// Case 003
		{
			lin: []byte(`Image: !Sub '${AWS::AccountId}.dkr.ecr.${AWS::Region}.amazonaws.com/splits-lite:bc7891268e44f62e0aebbe339c0850b61d52c417' # comment`),
			tag: []byte(`v3.5.0-bc789126`),
			out: []byte(`Image: !Sub '${AWS::AccountId}.dkr.ecr.${AWS::Region}.amazonaws.com/splits-lite:v3.5.0-bc789126' # comment`),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var err error

			var out []byte
			{
				out, err = repIma(tc.lin, tc.tag)
				if err != nil {
					t.Fatal("expected", nil, "got", err)
				}
			}

			if dif := cmp.Diff(tc.out, out); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
