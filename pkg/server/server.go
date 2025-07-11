package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/0xSplits/kayron/pkg/runtime"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twitchtv/twirp"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Int are the Twirp specific interceptors wrapping the endpoint handlers.
	Int []twirp.Interceptor
	// Lis is the main HTTP listener bound to some configured host and port.
	Lis net.Listener
	// Log is the structured logger passed down the stack.
	Log logger.Interface
	// Mid are the protocol specific transport layer middlewares executed before
	// any RPC handler.
	Mid []mux.MiddlewareFunc
}

type Server struct {
	lis net.Listener
	log logger.Interface
	srv *http.Server
}

func New(c Config) *Server {
	if c.Lis == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Lis must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	var rtr *mux.Router
	{
		rtr = mux.NewRouter()
	}

	{
		rtr.Use(c.Mid...)
	}

	// Add a simple health check response to the root.
	{
		rtr.NewRoute().Methods("GET").Path("/").HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
			wri.Header().Set("Content-Type", "text/plain")
			wri.WriteHeader(http.StatusOK)
			_, _ = wri.Write(linBrk([]byte("OK")))
		})
	}

	// Add a simple version response for the runtime.
	{
		rtr.NewRoute().Methods("GET").Path("/version").HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
			wri.Header().Set("Content-Type", "application/json")
			wri.WriteHeader(http.StatusOK)
			_, _ = wri.Write(linBrk(runtime.JSON()))
		})
	}

	// Add the metrics endpoint in Prometehus format.
	{
		rtr.NewRoute().Methods("GET").Path("/metrics").Handler(promhttp.Handler())
	}

	return &Server{
		lis: c.Lis,
		log: c.Log,
		srv: &http.Server{
			Handler: rtr,
		},
	}
}

func (s *Server) Daemon() {
	s.log.Log(
		"level", "info",
		"message", "server is accepting calls",
		"address", s.lis.Addr().String(),
	)

	{
		err := s.srv.Serve(s.lis)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}
}

func linBrk(byt []byte) []byte {
	return append(byt, []byte("\n")...)
}
