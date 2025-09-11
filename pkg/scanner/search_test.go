package scanner

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Scanner_Search(t *testing.T) {
	testCases := []struct {
		key string
	}{
		// Case 000
		{
			key: "  Service:",
		},
		// Case 001
		{
			key: "  TaskDefinition:",
		},
		// Case 002
		{
			key: "Resources:",
		},
		// Case 003
		{
			key: "      ServiceRegistries:",
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

			var sca *Scanner
			{
				sca = New(Config{
					Inp: inp,
				})
			}

			var res []byte
			{
				res = sca.Search([]byte(tc.key))
			}

			if dif := cmp.Diff(bytes.TrimSpace(out), bytes.TrimSpace(res)); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
