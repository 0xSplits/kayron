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
		sub []byte
	}{
		// Case 000
		{
			pre: "  Service:",
			sub: nil,
		},
		// Case 001
		{
			pre: "  TaskDefinition:",
			sub: nil,
		},
		// Case 002
		{
			pre: "Resources:",
			sub: nil,
		},
		// Case 003, real production example
		{
			pre: "      ServiceRegistries:",
			sub: nil,
		},
		// Case 004
		{
			pre: "          Image:",
			sub: nil,
		},
		// Case 005
		{
			pre: "          Image:",
			sub: []byte("          Image: registry/image:tag"),
		},
		// Case 006
		{
			pre: "      ContainerDefinitions:",
			sub: []byte("      ContainerDefinitions:\n        Foo: 1\n        Bar: 2"),
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
				res = sca.Delete([]byte(tc.pre), tc.sub...).Bytes()
			}

			if dif := cmp.Diff(bytes.TrimSpace(out), bytes.TrimSpace(res)); dif != "" {
				t.Fatalf("-expected +actual:\n%s", dif)
			}
		})
	}
}
