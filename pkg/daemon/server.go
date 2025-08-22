package daemon

import (
	"net"

	"github.com/0xSplits/kayron/pkg/server"
	"github.com/0xSplits/kayron/pkg/server/middleware/cors"
	"github.com/gorilla/mux"
	"github.com/xh3b4sd/tracer"
)

func (d *Daemon) Server() *server.Server {
	var err error

	var lis net.Listener
	{
		lis, err = net.Listen("tcp", net.JoinHostPort(d.env.HttpHost, d.env.HttpPort))
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	var ser *server.Server
	{
		ser = server.New(server.Config{
			Lis: lis,
			Log: d.log,
			Mid: []mux.MiddlewareFunc{
				cors.New(cors.Config{Log: d.log}).Handler,
			},
		})
	}

	return ser
}
