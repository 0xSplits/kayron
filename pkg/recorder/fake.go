package recorder

import "go.opentelemetry.io/otel/attribute"

type FakeConfig struct {
	Lab map[string][]string
}

type Fake struct {
	// lab is the set of labels to whitelist.
	lab map[string][]string
	// rec contains the information recorded during a single test run.
	rec *Recorded
}

type Recorded struct {
	Val []float64
	Lab []map[string]string
}

func NewFake(c FakeConfig) *Fake {
	return &Fake{
		lab: c.Lab,
		rec: &Recorded{},
	}
}

func (f *Fake) Labels() map[string][]string {
	return f.lab
}

func (f *Fake) Record(val float64, lab ...attribute.KeyValue) {
	{
		f.rec.Val = append(f.rec.Val, val)
	}

	m := map[string]string{}
	for _, x := range lab {
		m[string(x.Key)] = x.Value.AsString()
	}

	{
		f.rec.Lab = append(f.rec.Lab, m)
	}
}

func (f *Fake) Recorded() *Recorded {
	return f.rec
}
