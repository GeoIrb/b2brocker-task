package httpserver

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/valyala/fasthttp"
)

type service interface {
	Handler(ctx context.Context) (uuid string, err error)
}

type handlerServer struct {
	src             service
	transport       HandlerTransport
	errorProcessing errorProcessing

	logger log.Logger
}

func (s *handlerServer) ServeHTTP(ctx *fasthttp.RequestCtx) {
	level.Info(s.logger).Log("msg", "http request", "host", ctx.Request.Host())

	var (
		uuid string
		err  error
	)
	if uuid, err = s.src.Handler(context.Background()); err != nil {
		level.Error(s.logger).Log("msg", "proxy handler", "err", err)
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}

	if err = s.transport.Encode(&ctx.Response, uuid); err != nil {
		level.Error(s.logger).Log("msg", "encode responce", "err", err)
		s.errorProcessing(&ctx.Response, err, http.StatusInternalServerError)
		return
	}
}

// NewHandlerServer ...
func NewHandlerServer(
	src service,
	transport HandlerTransport,
	errorProcessing errorProcessing,

	logger log.Logger,
) fasthttp.RequestHandler {
	srv := &handlerServer{
		src:             src,
		transport:       transport,
		errorProcessing: errorProcessing,
		logger:          logger,
	}

	return srv.ServeHTTP
}
