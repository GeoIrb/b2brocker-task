package httprouter

import (
	"net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/go-kit/kit/log"
)

const (
	handlerHTTPMethod = http.MethodGet
	handlerURI        = "/handler"
)

// New return http server for proxy service
func New(svc service, logger log.Logger) *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.Handle(handlerHTTPMethod, handlerURI, NewHandlerServer(svc, NewHandlerTransport(), ErrorProcessing, logger))
	return router
}
