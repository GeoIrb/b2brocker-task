package proxy

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

// Proxy service
type Proxy struct {
	mutex   sync.Mutex
	request map[string]chan struct{}

	sendNext func(ctx context.Context, requestUUID string) error
}

// Handler to requests
func (p *Proxy) Handler(ctx context.Context) (err error) {
	uuid := uuid.New().String()
	ch := make(chan struct{})

	p.mutex.Lock()
	p.request[uuid] = ch
	p.mutex.Unlock()

	if err = p.sendNext(ctx, uuid); err != nil {
		p.mutex.Lock()
		delete(p.request, uuid)
		p.mutex.Unlock()
		return
	}
	<-ch
	return
}

// Compliter complites requests
func (p *Proxy) Compliter(requestUUID string) {
	p.mutex.Lock()
	if ch, isExist := p.request[requestUUID]; isExist {
		close(ch)
		delete(p.request, requestUUID)
	}
	p.mutex.Unlock()
}

// NewProxy service
func NewProxy(
	sendNext func(ctx context.Context, requestUUID string) error,
) *Proxy {
	return &Proxy{
		request:  make(map[string]chan struct{}),
		sendNext: sendNext,
	}
}
