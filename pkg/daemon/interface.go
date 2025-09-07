package daemon

import (
	"github.com/0xSplits/kayron/pkg/server"
	"github.com/0xSplits/workit/worker/combined"
)

type Interface interface {
	Server() *server.Server
	Worker() *combined.Worker
}
