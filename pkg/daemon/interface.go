package daemon

import (
	"github.com/0xSplits/kayron/pkg/server"
	"github.com/0xSplits/workit/worker/sequence"
)

type Interface interface {
	Server() *server.Server
	Worker() *sequence.Worker
}
