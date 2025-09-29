package webhook

import "time"

type Commit struct {
	Hash string
	Time time.Time
}

func (c Commit) Empty() bool {
	return c.Hash == "" && c.Time.IsZero()
}
