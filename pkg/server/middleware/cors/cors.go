package cors

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	Log logger.Interface
}

type Middleware struct {
	crs *cors.Cors
	log logger.Interface
}

func New(c Config) *Middleware {
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}

	return &Middleware{
		log: c.Log,
		crs: cors.New(cors.Options{
			ExposedHeaders: []string{"*"},
			AllowedHeaders: []string{"*"},
		}),
	}
}

func (m *Middleware) Handler(h http.Handler) http.Handler {
	return m.crs.Handler(h)
}
