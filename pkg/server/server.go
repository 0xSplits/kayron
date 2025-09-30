package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/0xSplits/kayron/pkg/runtime"
	"github.com/cbrgm/githubevents/v2/githubevents"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Lis is the main HTTP listener bound to some configured host and port.
	Lis net.Listener
	// Log is the structured logger passed down the stack.
	Log logger.Interface
	// Psh is the webhook event handler for push events. This event handler allows
	// us to cache real time data about commit hashes pushed across the entire
	// Github organization.
	Psh *githubevents.EventHandler
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
	if c.Psh == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Psh must not be empty", c)))
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

	// Add the metrics endpoint in Prometehus format.
	{
		rtr.NewRoute().Methods("GET").Path("/metrics").Handler(promhttp.Handler())
	}

	// Add a simple version response for the runtime.
	{
		rtr.NewRoute().Methods("GET").Path("/version").HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
			wri.Header().Set("Content-Type", "application/json")
			wri.WriteHeader(http.StatusOK)
			_, _ = wri.Write(linBrk(runtime.JSON()))
		})
	}

	// Add a webhook handler to receive near real time push events from Github.
	// Note that we should only accept the POST method here as per the Github
	// documentation.
	{
		rtr.NewRoute().Methods("POST").Path("/webhook/push").HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
			if err := c.Psh.HandleEventRequest(req); err != nil {
				http.Error(wri, err.Error(), http.StatusBadRequest)
			} else {
				wri.WriteHeader(http.StatusOK)
			}
		})
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
