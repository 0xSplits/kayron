package daemon

import (
	"github.com/0xSplits/kayron/pkg/server"
	"github.com/0xSplits/workit/worker"
)

type Interface interface {
	Server() *server.Server
	Worker() *worker.Worker
}
