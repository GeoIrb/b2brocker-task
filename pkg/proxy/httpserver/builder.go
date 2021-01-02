package httpserver

import (
	"net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/go-kit/kit/log"
)

const (
	handlerHTTPMethod = http.MethodGet
	handlerURI        = "/handler"
)

// NewServer return http server
func NewServer(svc service, logger log.Logger) *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.Handle(handlerHTTPMethod, handlerURI, NewHandlerServer(svc, NewHandlerTransport(), ErrorProcessing, logger))
	return router
}
