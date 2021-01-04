package mqserver

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"b2broker-task/pkg/mq"
)

type service interface {
	Handler(ctx context.Context, requestUUID string) (err error)
}

type handlerServer struct {
	srv             service
	transport       handlerTransport
	publishFunction mq.PublishFunction

	logger log.Logger
}

func (s *handlerServer) ServeMQ(ctx context.Context, data []byte) {
	level.Info(s.logger).Log("msg", "service message from mq", "message", string(data))

	uuid, err := s.transport.Decode(data)
	if err != nil {
		level.Error(s.logger).Log("msg", "service decode message from mq", "err", err)
		return
	}

	level.Info(s.logger).Log("msg", "service message from mq", "uuid", uuid)
	if err := s.srv.Handler(ctx, uuid); err != nil {
		level.Error(s.logger).Log("msg", "service handle message from mq", "err", err)
		return
	}

	message, err := s.transport.Encode(uuid)
	if err != nil {
		level.Error(s.logger).Log("msg", "service encode handle message from mq", "err", err)
		return
	}

	if err = s.publishFunction(ctx, message); err != nil {
		level.Error(s.logger).Log("msg", "service publish message to mq", "err", err)
		return
	}
}

// NewHandlerServer ...
func NewHandlerServer(
	srv service,
	transport handlerTransport,
	publishFunction mq.PublishFunction,
	logger log.Logger,
) mq.Handler {
	s := &handlerServer{
		srv:             srv,
		transport:       transport,
		publishFunction: publishFunction,
		logger:          logger,
	}

	return s.ServeMQ
}
