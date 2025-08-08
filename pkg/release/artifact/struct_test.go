package artifact

import (
	"testing"

	"github.com/0xSplits/kayron/pkg/release/artifact/condition"
	"github.com/0xSplits/kayron/pkg/release/artifact/reference"
	"github.com/0xSplits/kayron/pkg/release/artifact/scheduler"
	"github.com/google/go-cmp/cmp"
)

func Test_Release_Artifact_Merge(t *testing.T) {
	var act Struct

	{
		act = act.Merge(Struct{
			Condition: condition.Struct{
				Success: true,
			},
		})
	}

	{
		act = act.Merge(Struct{
			Reference: reference.Struct{
				Desired: "desired",
			},
		})
	}

	{
		act = act.Merge(Struct{
			Scheduler: scheduler.Struct{
				Current: "current",
			},
		})
	}

	var exp Struct
	{
		exp = Struct{
			Condition: condition.Struct{
				Success: true,
			},
			Reference: reference.Struct{
				Desired: "desired",
			},
			Scheduler: scheduler.Struct{
				Current: "current",
			},
		}
	}

	if dif := cmp.Diff(exp, act); dif != "" {
		t.Fatalf("-expected +actual:\n%s", dif)
	}
}
