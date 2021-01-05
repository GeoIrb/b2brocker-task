package mqhandler

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"b2broker-task/pkg/mq"
)

type service interface {
	Compliter(ctx context.Context, requestUUID string)
}

type handlerServer struct {
	srv       service
	transport HandlerTransport

	logger log.Logger
}

func (s *handlerServer) ServeMQ(ctx context.Context, data []byte) {
	level.Info(s.logger).Log("msg", "service message from mq", "message", string(data))

	uuid, err := s.transport.Decode(data)
	if err != nil {
		level.Error(s.logger).Log("msg", "service decode message from mq", "err", err)
		return
	}

	level.Error(s.logger).Log("msg", "proxy", "uuid", uuid)
	s.srv.Compliter(ctx, uuid)
}

// NewHandlerServer ...
func NewHandlerServer(
	src service,
	transport HandlerTransport,
	logger log.Logger,
) mq.Handler {
	srv := &handlerServer{
		srv:       src,
		transport: transport,
		logger:    logger,
	}

	return srv.ServeMQ
}
