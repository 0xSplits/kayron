package middleware

import "net/http"

type Interface interface {
	Handler(http.Handler) http.Handler
}
