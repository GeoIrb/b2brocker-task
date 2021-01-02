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

	sendNext func(requestUUID string) error
}

// Handler to requests
func (p *Proxy) Handler(ctx context.Context) (requestUUID string, err error) {
	requestUUID = uuid.New().String()
	ch := make(chan struct{})

	p.mutex.Lock()
	p.request[requestUUID] = ch
	p.mutex.Unlock()

	if err = p.sendNext(requestUUID); err != nil {
		p.mutex.Lock()
		delete(p.request, requestUUID)
		p.mutex.Unlock()
		return
	}
	<-ch
	return
}

// Compliter complites requests
func (p *Proxy) Compliter(ctx context.Context, requestUUID string) {
	p.mutex.Lock()
	if ch, isExist := p.request[requestUUID]; isExist {
		close(ch)
		delete(p.request, requestUUID)
	}
	p.mutex.Unlock()
}

// NewProxy service
func NewProxy(
	sendNext func(requestUUID string) error,
) *Proxy {
	return &Proxy{
		request:  make(map[string]chan struct{}),
		sendNext: sendNext,
	}
}
