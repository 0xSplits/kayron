package preview

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Operator_Infrastructure_Preview_Render(t *testing.T) {
	testCases := []struct {
		bra []string
	}{
		// Case 000
		{
			bra: []string{
				"fancy-feature-branch",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var err error

			var inp []byte
			{
				inp, err = os.ReadFile(fmt.Sprintf("./testdata/%03d/inp.yaml.golden", i))
				if err != nil {
					t.Fatal("expected", nil, "got", err)
				}
			}

			var out []byte
			{
				out, err = os.ReadFile(fmt.Sprintf("./testdata/%03d/out.yaml.golden", i))
				if err != nil {
					t.Fatal("expected", nil, "got", err)
				}
			}

			var pre *Preview
			{
				pre = New(Config{
					Inp: inp,
				})
			}

			var res []byte
			{
				res = pre.Render(tc.bra)
				if err != nil {
					t.Fatal("expected", nil, "got", err)
				}
			}

			if dif := cmp.Diff(bytes.TrimSpace(out), bytes.TrimSpace(res)); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
