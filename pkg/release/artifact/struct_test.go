package artifact

import (
	"testing"

	"github.com/0xSplits/kayron/pkg/release/artifact/condition"
	"github.com/0xSplits/kayron/pkg/release/artifact/reference"
	"github.com/0xSplits/kayron/pkg/release/artifact/scheduler"
	"github.com/google/go-cmp/cmp"
)

// Test_Release_Artifact_Merge_forward ensures that non-zero values can always
// overwrite zero values.
func Test_Release_Artifact_Merge_forward(t *testing.T) {
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
			Condition: condition.Struct{
				Trigger: true,
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
				Trigger: true,
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

// Test_Release_Artifact_Merge_no_overwrite ensures that non-zero values are
// never overwritten.
func Test_Release_Artifact_Merge_no_overwrite(t *testing.T) {
	var act Struct
	{
		act = Struct{
			Condition: condition.Struct{
				Success: true,
				Trigger: true,
			},
			Reference: reference.Struct{
				Desired: "desired",
			},
			Scheduler: scheduler.Struct{
				Current: "current",
			},
		}
	}

	{
		act = act.Merge(Struct{
			Condition: condition.Struct{
				Success: false,
			},
		})
	}

	{
		act = act.Merge(Struct{
			Condition: condition.Struct{
				Trigger: false,
			},
		})
	}

	{
		act = act.Merge(Struct{
			Reference: reference.Struct{
				Desired: "no bueno",
			},
		})
	}

	{
		act = act.Merge(Struct{
			Scheduler: scheduler.Struct{
				Current: "foobar",
			},
		})
	}

	var exp Struct
	{
		exp = Struct{
			Condition: condition.Struct{
				Success: true,
				Trigger: true,
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
