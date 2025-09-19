package scanner

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Scanner_Append(t *testing.T) {
	testCases := []struct {
		pre string
		suf string
	}{
		// Case 000
		{
			pre: "    Value:",
			suf: "-1234",
		},
		// Case 001
		{
			pre: "    Value:",
			suf: ".0xFa73",
		},
		// Case 002
		{
			pre: "      ServiceName:",
			suf: "-1d0fd508",
		},
		// Case 003
		{
			pre: "      TaskDefinition:",
			suf: "-e3eae11",
		},
		// Case 004
		{
			pre: "      TaskDefinition:",
			suf: "-XHEKSOUDL",
		},
		// Case 005
		{
			pre: "  Service:",
			suf: "FancyFeatureBranch",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			var err error

			var inp []byte
			{
				inp, err = os.ReadFile(fmt.Sprintf("./testdata/append/%03d/inp.yaml.golden", i))
				if err != nil {
					t.Fatal("expected", nil, "got", err)
				}
			}

			var out []byte
			{
				out, err = os.ReadFile(fmt.Sprintf("./testdata/append/%03d/out.yaml.golden", i))
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
				res = sca.Append([]byte(tc.pre), []byte(tc.suf)).Bytes()
			}

			if dif := cmp.Diff(bytes.TrimSpace(out), bytes.TrimSpace(res)); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
