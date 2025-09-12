package scanner

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Scanner_Delete(t *testing.T) {
	testCases := []struct {
		pre string
	}{
		// Case 000
		{
			pre: "  Service:",
		},
		// Case 001
		{
			pre: "  TaskDefinition:",
		},
		// Case 002
		{
			pre: "Resources:",
		},
		// Case 003
		{
			pre: "      ServiceRegistries:",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var err error

			var inp []byte
			{
				inp, err = os.ReadFile(fmt.Sprintf("./testdata/delete/%03d/inp.yaml.golden", i))
				if err != nil {
					t.Fatal("expected", nil, "got", err)
				}
			}

			var out []byte
			{
				out, err = os.ReadFile(fmt.Sprintf("./testdata/delete/%03d/out.yaml.golden", i))
				if err != nil {
					t.Fatal("expected", nil, "got", err)
				}
			}

			var sca *Scanner
			{
				sca = New(Config{
					Inp: inp,
				})
			}

			var res []byte
			{
				res = sca.Delete([]byte(tc.pre)).Bytes()
			}

			if dif := cmp.Diff(bytes.TrimSpace(out), bytes.TrimSpace(res)); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
